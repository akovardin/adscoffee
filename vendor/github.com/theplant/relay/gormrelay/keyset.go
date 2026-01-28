package gormrelay

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/theplant/relay"
	"github.com/theplant/relay/cursor"
)

func createWhereExpr(getColumn func(fieldName string) (clause.Column, error), orderBys []relay.OrderBy, keyset map[string]any, reverse bool) (clause.Expression, error) {
	ors := make([]clause.Expression, 0, len(orderBys))
	eqs := make([]clause.Expression, 0, len(orderBys))
	for i, orderBy := range orderBys {
		v, ok := keyset[orderBy.Field]
		if !ok {
			return nil, errors.Errorf("missing field %q in keyset", orderBy.Field)
		}

		column, err := getColumn(orderBy.Field)
		if err != nil {
			return nil, err
		}

		desc := orderBy.Desc
		if reverse {
			desc = !desc
		}

		var expr clause.Expression
		if desc {
			expr = clause.Lt{Column: column, Value: v}
		} else {
			expr = clause.Gt{Column: column, Value: v}
		}

		ands := make([]clause.Expression, len(eqs)+1)
		copy(ands, eqs)
		ands[len(eqs)] = expr
		ors = append(ors, clause.And(ands...))

		if i < len(orderBys)-1 {
			eqs = append(eqs, clause.Eq{Column: column, Value: v})
		}
	}
	return clause.And(clause.Or(ors...)), nil
}

// Example:
// db.Clauses(
//
//	 	// This is for `Where`, so we cant use `Where(clause.And(clause.Or(...),clause.Or(...)))`
//		clause.And(
//			clause.Or( // after
//				clause.And(
//					clause.Gt{Column: "age", Value: 85}, // ASC
//				),
//				clause.And(
//					clause.Eq{Column: "age", Value: 85},
//					clause.Lt{Column: "name", Value: "name15"}, // DESC
//				),
//			),
//		),
//		clause.And(
//			clause.Or( // before
//				clause.And(
//					clause.Lt{Column: "age", Value: 88},
//				),
//				clause.And(
//					clause.Eq{Column: "age", Value: 88},
//					clause.Gt{Column: "name", Value: "name12"},
//				),
//			),
//		),
//		clause.OrderBy{
//			Columns: []clause.OrderByColumn{
//				{Column: clause.Column{Name: "age"}, Desc: false},
//				{Column: clause.Column{Name: "name"}, Desc: true},
//			},
//		},
//		clause.Limit{Limit: &limit},
//
// )
func scopeKeyset(computedColumns map[string]clause.Column, after, before *map[string]any, orderBys []relay.OrderBy, limit int, fromEnd bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if db.Statement.Model == nil {
			_ = db.AddError(errors.New("model is nil"))
			return db
		}

		s, err := parseSchema(db, db.Statement.Model)
		if err != nil {
			_ = db.AddError(err)
			return db
		}

		getColumn := func(fieldName string) (clause.Column, error) {
			if column, ok := computedColumns[fieldName]; ok {
				return column, nil
			}
			field, ok := s.FieldsByName[fieldName]
			if !ok {
				return clause.Column{}, errors.Errorf("missing field %q in schema", fieldName)
			}
			return clause.Column{Table: clause.CurrentTable, Name: field.DBName}, nil
		}

		var exprs []clause.Expression

		if after != nil {
			expr, err := createWhereExpr(getColumn, orderBys, *after, false)
			if err != nil {
				_ = db.AddError(err)
				return db
			}
			exprs = append(exprs, expr)
		}

		if before != nil {
			expr, err := createWhereExpr(getColumn, orderBys, *before, true)
			if err != nil {
				_ = db.AddError(err)
				return db
			}
			exprs = append(exprs, expr)
		}

		computedColumns = lo.MapEntries(computedColumns, func(field string, column clause.Column) (string, clause.Column) {
			column.Alias = ComputedFieldToColumnAlias(field)
			return field, column
		})

		if len(orderBys) > 0 {
			orderByColumns := make([]clause.OrderByColumn, 0, len(orderBys))
			for _, orderBy := range orderBys {
				column, ok := computedColumns[orderBy.Field]
				if ok {
					column = clause.Column{Name: column.Alias, Raw: true}
				} else {
					column, err = getColumn(orderBy.Field)
					if err != nil {
						_ = db.AddError(err)
						return db
					}
				}

				desc := orderBy.Desc
				if fromEnd {
					desc = !desc
				}
				orderByColumns = append(orderByColumns, clause.OrderByColumn{
					Column: column,
					Desc:   desc,
				})
			}
			exprs = append(exprs, clause.OrderBy{Columns: orderByColumns})
		}

		if limit > 0 {
			exprs = append(exprs, clause.Limit{Limit: &limit})
		} else {
			_ = db.AddError(errors.New("limit must be greater than 0"))
		}
		if len(computedColumns) > 0 {
			db = db.Scopes(AppendSelect(maps.Values(computedColumns)...))
		}
		return db.Clauses(exprs...)
	}
}

type KeysetFinder[T any] struct {
	db   *gorm.DB
	opts Options[T]
}

func NewKeysetFinder[T any](db *gorm.DB, opts ...Option[T]) *KeysetFinder[T] {
	options := Options[T]{}
	for _, opt := range opts {
		opt(&options)
	}
	return &KeysetFinder[T]{db: db, opts: options}
}

func (a *KeysetFinder[T]) Find(ctx context.Context, after, before *map[string]any, orderBys []relay.OrderBy, limit int, fromEnd bool) ([]cursor.Node[T], error) {
	if limit == 0 {
		return []cursor.Node[T]{}, nil
	}

	db := a.db.WithContext(ctx)

	basedOnModel, err := shouldBasedOnModel[T](db)
	if err != nil {
		return nil, err
	}

	if !basedOnModel {
		db = applyModel[T](db)
	}

	if a.opts.Computed != nil {
		if err := a.opts.Computed.Validate(); err != nil {
			return nil, err
		}

		dest, toCursorNodes, err := a.opts.Computed.ForScan(db)
		if err != nil {
			return nil, err
		}

		computedResults, err := splitComputedScan(
			a.opts.Computed.Columns,
			db.Scopes(scopeKeyset(a.opts.Computed.Columns, after, before, orderBys, limit, fromEnd)),
			dest,
		)
		if err != nil {
			return nil, err
		}

		nodes := toCursorNodes(computedResults)
		if fromEnd {
			lo.Reverse(nodes)
		}
		return nodes, nil
	}

	if basedOnModel {
		modelType := reflect.TypeOf(db.Statement.Model)
		sliceType := reflect.SliceOf(modelType)
		nodesVal := reflect.New(sliceType).Elem()

		err := db.Scopes(scopeKeyset(nil, after, before, orderBys, limit, fromEnd)).Find(nodesVal.Addr().Interface()).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find records based on model")
		}

		nodes := make([]T, nodesVal.Len())
		for i := 0; i < nodesVal.Len(); i++ {
			nodes[i] = nodesVal.Index(i).Interface().(T)
		}

		if fromEnd {
			lo.Reverse(nodes)
		}
		return lo.Map(nodes, func(node T, _ int) cursor.Node[T] {
			return &cursor.SelfNode[T]{Node: node}
		}), nil
	}

	var nodes []T
	err = db.Scopes(scopeKeyset(nil, after, before, orderBys, limit, fromEnd)).Find(&nodes).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find records with keyset pagination")
	}
	if fromEnd {
		lo.Reverse(nodes)
	}
	return lo.Map(nodes, func(node T, _ int) cursor.Node[T] {
		return &cursor.SelfNode[T]{Node: node}
	}), nil
}

func (a *KeysetFinder[T]) Count(ctx context.Context) (int, error) {
	db := a.db.WithContext(ctx)

	basedOnModel, err := shouldBasedOnModel[T](db)
	if err != nil {
		return 0, err
	}

	if !basedOnModel {
		db = applyModel[T](db)
	}

	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return 0, errors.Wrap(err, "failed to count total records")
	}
	return int(totalCount), nil
}

var _ cursor.KeysetFinder[any] = &KeysetFinder[any]{}

func NewKeysetAdapter[T any](db *gorm.DB, opts ...Option[T]) relay.ApplyCursorsFunc[T] {
	return cursor.NewKeysetAdapter(NewKeysetFinder(db, opts...))
}
