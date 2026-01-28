package relay

import (
	"context"

	"github.com/samber/lo"
)

// PaginationMiddleware is a wrapper for Pagination (middleware pattern)
type PaginationMiddleware[T any] func(next Pagination[T]) Pagination[T]

// EnsureLimits ensures that the limit is within the range 0 -> maxLimit and uses defaultLimit if limit is not set or is negative
// This method introduced a breaking change in version 0.4.0, intentionally swapping the order of parameters to strongly indicate the breaking change.
// https://github.com/theplant/relay/compare/genx?expand=1#diff-02f50901140d6057da6310a106670552aa766a093efbc2200fb34c099b762131R14
func EnsureLimits[T any](defaultLimit, maxLimit int) PaginationMiddleware[T] {
	if defaultLimit < 0 {
		panic("defaultLimit cannot be negative")
	}
	if maxLimit < defaultLimit {
		panic("maxLimit must be greater than or equal to defaultLimit")
	}
	return func(next Pagination[T]) Pagination[T] {
		return PaginationFunc[T](func(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error) {
			if req.First != nil {
				if *req.First > maxLimit {
					req.First = &maxLimit
				}
				if *req.First < 0 {
					req.First = &defaultLimit
				}
			}
			if req.Last != nil {
				if *req.Last > maxLimit {
					req.Last = &maxLimit
				}
				if *req.Last < 0 {
					req.Last = &defaultLimit
				}
			}
			if req.First == nil && req.Last == nil {
				if req.After == nil && req.Before != nil {
					req.Last = &defaultLimit
				} else {
					req.First = &defaultLimit
				}
			}
			return next.Paginate(ctx, req)
		})
	}
}

func EnsurePrimaryOrderBy[T any](primaryOrderBys ...OrderBy) PaginationMiddleware[T] {
	return func(next Pagination[T]) Pagination[T] {
		return PaginationFunc[T](func(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error) {
			req.OrderBys = AppendPrimaryOrderBy(req.OrderBys, primaryOrderBys...)
			return next.Paginate(ctx, req)
		})
	}
}

func AppendPrimaryOrderBy(orderBys []OrderBy, primaryOrderBys ...OrderBy) []OrderBy {
	if len(primaryOrderBys) == 0 {
		return orderBys
	}
	orderByFields := lo.SliceToMap(orderBys, func(orderBy OrderBy) (string, bool) {
		return orderBy.Field, true
	})
	// If there are fields in primaryOrderBys that are not in orderBys, add them to orderBys
	for _, primaryOrderBy := range primaryOrderBys {
		if _, ok := orderByFields[primaryOrderBy.Field]; !ok {
			orderBys = append(orderBys, primaryOrderBy)
		}
	}
	return orderBys
}

// CursorMiddleware is a wrapper for ApplyCursorsFunc (middleware pattern)
type CursorMiddleware[T any] func(next ApplyCursorsFunc[T]) ApplyCursorsFunc[T]

type ctxCursorMiddlewares struct{}

func CursorMiddlewaresFromContext[T any](ctx context.Context) []CursorMiddleware[T] {
	middlewares, _ := ctx.Value(ctxCursorMiddlewares{}).([]CursorMiddleware[T])
	return middlewares
}

func AppendCursorMiddleware[T any](cursorMiddlewares ...CursorMiddleware[T]) PaginationMiddleware[T] {
	return func(next Pagination[T]) Pagination[T] {
		return PaginationFunc[T](func(ctx context.Context, req *PaginateRequest[T]) (*Connection[T], error) {
			if len(cursorMiddlewares) > 0 {
				cursorMiddlewares := append(CursorMiddlewaresFromContext[T](ctx), cursorMiddlewares...)
				ctx = context.WithValue(ctx, ctxCursorMiddlewares{}, cursorMiddlewares)
			}
			return next.Paginate(ctx, req)
		})
	}
}

func chainCursorMiddlewares[T any](mws []CursorMiddleware[T]) CursorMiddleware[T] {
	return func(next ApplyCursorsFunc[T]) ApplyCursorsFunc[T] {
		for i := len(mws); i > 0; i-- {
			next = mws[i-1](next)
		}
		return next
	}
}

func chainPaginationMiddlewares[T any](mws []PaginationMiddleware[T]) PaginationMiddleware[T] {
	return func(next Pagination[T]) Pagination[T] {
		for i := len(mws); i > 0; i-- {
			next = mws[i-1](next)
		}
		return next
	}
}
