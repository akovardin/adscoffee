package cursor

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/theplant/relay"
)

type OffsetFinder[T any] interface {
	Find(ctx context.Context, orderBys []relay.OrderBy, skip, limit int) ([]Node[T], error)
	Count(ctx context.Context) (int, error)
}

type OffsetFinderFunc[T any] func(ctx context.Context, orderBys []relay.OrderBy, skip, limit int) ([]Node[T], error)

func (f OffsetFinderFunc[T]) Find(ctx context.Context, orderBys []relay.OrderBy, skip, limit int) ([]Node[T], error) {
	return f(ctx, orderBys, skip, limit)
}

// NewOffsetAdapter creates a relay.ApplyCursorsFunc from an OffsetFinder.
// If you want to use `last!=nil&&before==nil`, you can't skip totalCount.
func NewOffsetAdapter[T any](finder OffsetFinder[T]) relay.ApplyCursorsFunc[T] {
	return func(ctx context.Context, req *relay.ApplyCursorsRequest) (*relay.ApplyCursorsResponse[T], error) {
		after, before, err := decodeOffsetCursors(req.After, req.Before)
		if err != nil {
			return nil, err
		}

		var skipCount, skipFind bool
		{
			skip := relay.GetSkip(ctx)
			skipCount = skip.TotalCount
			skipFind = skip.Nodes && skip.Edges && skip.PageInfo
		}

		var totalCount *int
		if !skipCount {
			count, err := finder.Count(ctx)
			if err != nil {
				return nil, err
			}
			totalCount = &count
		}

		if skipFind {
			return &relay.ApplyCursorsResponse[T]{
				TotalCount: totalCount,
			}, nil
		}

		if req.FromEnd && before == nil {
			if totalCount == nil {
				return nil, errors.New("totalCount is required for pagination from end when before cursor is not provided")
			}
			before = totalCount
		}

		limit, skip := req.Limit, 0
		if after != nil {
			skip = *after + 1
		} else if before != nil {
			skip = *before - limit
		}
		if skip < 0 {
			skip = 0
		}
		if before != nil {
			rangeLen := *before - skip
			if rangeLen <= 0 {
				rangeLen = 0
			}
			if limit > rangeLen {
				limit = rangeLen
			}
			if req.FromEnd && limit < rangeLen {
				skip = *before - limit
			}
		}

		var edges []*relay.LazyEdge[T]
		if limit <= 0 || (totalCount != nil && (skip >= *totalCount || *totalCount <= 0)) {
			edges = make([]*relay.LazyEdge[T], 0)
		} else {
			nodes, err := finder.Find(ctx, req.OrderBys, skip, limit)
			if err != nil {
				return nil, err
			}
			edges = make([]*relay.LazyEdge[T], len(nodes))
			for i, node := range nodes {
				i := i
				edges[i] = &relay.LazyEdge[T]{
					Node: node.RelayNode(),
					Cursor: func(_ context.Context) (string, error) {
						return EncodeOffsetCursor(skip + i), nil
					},
				}
			}
		}

		rsp := &relay.ApplyCursorsResponse[T]{
			LazyEdges:  edges,
			TotalCount: totalCount,
		}

		if totalCount != nil {
			rsp.HasAfterOrPrevious = after != nil && *after < *totalCount
			rsp.HasBeforeOrNext = before != nil && *before < *totalCount
		} else {
			// If we don't have totalCount, it would be very costly to check whether after and before really exist,
			// So it is usually not worth it. Normally, checking that it is not nil is sufficient.
			rsp.HasAfterOrPrevious = after != nil
			rsp.HasBeforeOrNext = before != nil
		}

		return rsp, nil
	}
}

func EncodeOffsetCursor(offset int) string {
	return strconv.Itoa(offset)
}

func DecodeOffsetCursor(cursor string) (int, error) {
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid offset cursor %q: cannot convert to integer", cursor)
	}
	return offset, nil
}

func decodeOffsetCursors(after, before *string) (afterOffset, beforeOffset *int, err error) {
	if after != nil {
		offset, err := DecodeOffsetCursor(*after)
		if err != nil {
			return nil, nil, err
		}
		afterOffset = &offset
	}
	if before != nil {
		offset, err := DecodeOffsetCursor(*before)
		if err != nil {
			return nil, nil, err
		}
		beforeOffset = &offset
	}
	if afterOffset != nil && *afterOffset < 0 {
		return nil, nil, errors.New("invalid pagination: after cursor must be non-negative")
	}
	if beforeOffset != nil && *beforeOffset < 0 {
		return nil, nil, errors.New("invalid pagination: before cursor must be non-negative")
	}
	if afterOffset != nil && before != nil && *afterOffset >= *beforeOffset {
		return nil, nil, errors.New("invalid pagination: after cursor must be less than before cursor")
	}
	return afterOffset, beforeOffset, nil
}
