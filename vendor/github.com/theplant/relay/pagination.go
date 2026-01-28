package relay

import (
	"context"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type OrderBy struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}

type PaginateRequest[T any] struct {
	After    *string   `json:"after"`
	First    *int      `json:"first"`
	Before   *string   `json:"before"`
	Last     *int      `json:"last"`
	OrderBys []OrderBy `json:"orderBys"`
}

type Edge[T any] struct {
	Node   T      `json:"node"`
	Cursor string `json:"cursor"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

type Connection[T any] struct {
	Edges      []*Edge[T] `json:"edges,omitempty"`
	Nodes      []T        `json:"nodes,omitempty"`
	PageInfo   *PageInfo  `json:"pageInfo,omitempty"`
	TotalCount *int       `json:"totalCount,omitempty"`
}

type ApplyCursorsRequest struct {
	Before   *string
	After    *string
	OrderBys []OrderBy
	Limit    int
	FromEnd  bool
}

type LazyEdge[T any] struct {
	Node   T
	Cursor func(ctx context.Context) (string, error)
}

type ApplyCursorsResponse[T any] struct {
	LazyEdges          []*LazyEdge[T]
	TotalCount         *int
	HasBeforeOrNext    bool // `before` exists or it's next exists
	HasAfterOrPrevious bool // `after` exists or it's previous exists
}

// https://relay.dev/graphql/connections.htm#ApplyCursorsToEdges()
type ApplyCursorsFunc[T any] func(ctx context.Context, req *ApplyCursorsRequest) (*ApplyCursorsResponse[T], error)

// https://relay.dev/graphql/connections.htm#sec-Pagination-algorithm
// https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo.Fields
func paginate[T any](ctx context.Context, req *PaginateRequest[T], applyCursorsFunc ApplyCursorsFunc[T]) (*Connection[T], error) {
	if req.First == nil && req.Last == nil {
		return nil, errors.New("first or last must be set")
	}
	if req.First != nil && req.Last != nil {
		return nil, errors.New("first and last cannot be used together")
	}
	if req.First != nil && *req.First < 0 {
		return nil, errors.New("first must be a non-negative integer")
	}
	if req.Last != nil && *req.Last < 0 {
		return nil, errors.New("last must be a non-negative integer")
	}

	orderBys := req.OrderBys
	if len(orderBys) > 0 {
		dups := lo.FindDuplicatesBy(orderBys, func(item OrderBy) string {
			return item.Field
		})
		if len(dups) > 0 {
			return nil, errors.Errorf("duplicated order by fields %v", lo.Map(dups, func(item OrderBy, _ int) string {
				return item.Field
			}))
		}
	}

	skip := GetSkip(ctx)
	if skip.All() {
		return &Connection[T]{}, nil
	}

	var limit int
	if req.First != nil {
		limit = *req.First + 1
	} else {
		limit = *req.Last + 1
	}

	rsp, err := applyCursorsFunc(ctx, &ApplyCursorsRequest{
		Before:   req.Before,
		After:    req.After,
		OrderBys: orderBys,
		Limit:    limit,
		FromEnd:  req.Last != nil,
	})
	if err != nil {
		return nil, err
	}

	lazyEdges := rsp.LazyEdges

	processor := GetNodeProcessor[T](ctx)
	if processor != nil {
		for _, lazyEdge := range lazyEdges {
			node, err := processor(ctx, lazyEdge.Node)
			if err != nil {
				return nil, err
			}
			lazyEdge.Node = node
		}
	}

	var hasPreviousPage, hasNextPage bool

	if req.First != nil && len(lazyEdges) > *req.First {
		lazyEdges = lazyEdges[:*req.First]
		hasNextPage = true
	}
	if req.Before != nil && rsp.HasBeforeOrNext {
		hasNextPage = true
	}

	if req.Last != nil && len(lazyEdges) > *req.Last {
		lazyEdges = lazyEdges[len(lazyEdges)-*req.Last:]
		hasPreviousPage = true
	}
	if req.After != nil && rsp.HasAfterOrPrevious {
		hasPreviousPage = true
	}

	conn := &Connection[T]{}

	if !skip.Edges {
		edges := make([]*Edge[T], len(lazyEdges))
		for i, lazyEdge := range lazyEdges {
			cursor, err := lazyEdge.Cursor(ctx)
			if err != nil {
				return nil, err
			}
			edges[i] = &Edge[T]{Node: lazyEdge.Node, Cursor: cursor}
		}
		conn.Edges = edges
	}

	if !skip.Nodes {
		nodes := make([]T, len(lazyEdges))
		for i, lazyEdge := range lazyEdges {
			nodes[i] = lazyEdge.Node
		}
		conn.Nodes = nodes
	}

	if !skip.TotalCount {
		conn.TotalCount = rsp.TotalCount
	}

	if !skip.PageInfo {
		pageInfo := &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
		}
		if len(lazyEdges) > 0 {
			var startCursor, endCursor string
			if len(conn.Edges) > 0 {
				startCursor = conn.Edges[0].Cursor
				endCursor = conn.Edges[len(conn.Edges)-1].Cursor
			} else {
				startCursor, err = lazyEdges[0].Cursor(ctx)
				if err != nil {
					return nil, err
				}
				endIndex := len(lazyEdges) - 1
				if endIndex == 0 {
					endCursor = startCursor
				} else {
					endCursor, err = lazyEdges[endIndex].Cursor(ctx)
					if err != nil {
						return nil, err
					}
				}
			}
			pageInfo.StartCursor = &startCursor
			pageInfo.EndCursor = &endCursor
		}
		conn.PageInfo = pageInfo
	}

	return conn, nil
}

type Pagination[T any] interface {
	Paginate(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error)
}

type PaginationFunc[T any] func(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error)

func (f PaginationFunc[T]) Paginate(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error) {
	return f(ctx, req)
}

func New[T any](applyCursorsFunc ApplyCursorsFunc[T], middlewares ...PaginationMiddleware[T]) Pagination[T] {
	if applyCursorsFunc == nil {
		panic("applyCursorsFunc must be set")
	}

	var p Pagination[T] = PaginationFunc[T](func(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error) {
		cursorMiddlewares := CursorMiddlewaresFromContext[T](ctx)
		return paginate(ctx, req, chainCursorMiddlewares(cursorMiddlewares)(applyCursorsFunc))
	})
	return Wrap(p, middlewares...)
}

func Wrap[T any](p Pagination[T], middlewares ...PaginationMiddleware[T]) Pagination[T] {
	return chainPaginationMiddlewares(middlewares)(p)
}
