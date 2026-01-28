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

type OffsetFinder[T any] struct {
	db   *gorm.DB
	opts Options[T]
}

func NewOffsetFinder[T any](db *gorm.DB, opts ...Option[T]) *OffsetFinder[T] {
	options := Options[T]{}
	for _, opt := range opts {
		opt(&options)
	}
	return &OffsetFinder[T]{db: db, opts: options}
}

func (a *OffsetFinder[T]) Find(ctx context.Context, orderBys []relay.OrderBy, skip, limit int) ([]cursor.Node[T], error) {
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

	if skip > 0 {
		db = db.Offset(skip)
	}

	db = db.Limit(limit)

	var computedColumns map[string]clause.Column
	if a.opts.Computed != nil {
		if err := a.opts.Computed.Validate(); err != nil {
			return nil, err
		}

		computedColumns = lo.MapEntries(a.opts.Computed.Columns, func(field string, column clause.Column) (string, clause.Column) {
			column.Alias = ComputedFieldToColumnAlias(field)
			return field, column
		})
	}

	if len(orderBys) > 0 {
		s, err := parseSchema(db, db.Statement.Model)
		if err != nil {
			return nil, err
		}

		getColumn := func(fieldName string) (clause.Column, error) {
			field, ok := s.FieldsByName[fieldName]
			if !ok {
				return clause.Column{}, errors.Errorf("missing field %q in schema for OrderBy", fieldName)
			}
			return clause.Column{Table: clause.CurrentTable, Name: field.DBName}, nil
		}

		orderByColumns := make([]clause.OrderByColumn, 0, len(orderBys))
		for _, orderBy := range orderBys {
			column, ok := computedColumns[orderBy.Field]
			if ok {
				column = clause.Column{Name: column.Alias, Raw: true}
			} else {
				column, err = getColumn(orderBy.Field)
				if err != nil {
					return nil, err
				}
			}

			orderByColumns = append(orderByColumns, clause.OrderByColumn{
				Column: column,
				Desc:   orderBy.Desc,
			})
		}
		db = db.Order(clause.OrderBy{Columns: orderByColumns})
	}

	if len(computedColumns) > 0 {
		dest, toCursorNodes, err := a.opts.Computed.ForScan(db)
		if err != nil {
			return nil, err
		}

		computedResults, err := splitComputedScan(
			computedColumns,
			db.Scopes(AppendSelect(maps.Values(computedColumns)...)),
			dest,
		)
		if err != nil {
			return nil, err
		}

		return toCursorNodes(computedResults), nil
	}

	var nodes []T

	if basedOnModel {
		modelType := reflect.TypeOf(db.Statement.Model)
		sliceType := reflect.SliceOf(modelType)
		nodesVal := reflect.New(sliceType).Elem()

		err := db.Find(nodesVal.Addr().Interface()).Error
		if err != nil {
			return nil, errors.Wrap(err, "failed to find records based on model")
		}

		nodes = make([]T, nodesVal.Len())
		for i := 0; i < nodesVal.Len(); i++ {
			nodes[i] = nodesVal.Index(i).Interface().(T)
		}
	} else {
		if err := db.Find(&nodes).Error; err != nil {
			return nil, errors.Wrap(err, "failed to find records with offset pagination")
		}
	}

	return lo.Map(nodes, func(node T, _ int) cursor.Node[T] {
		return &cursor.SelfNode[T]{Node: node}
	}), nil
}

func (a *OffsetFinder[T]) Count(ctx context.Context) (int, error) {
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

var _ cursor.OffsetFinder[any] = &OffsetFinder[any]{}

func NewOffsetAdapter[T any](db *gorm.DB, opts ...Option[T]) relay.ApplyCursorsFunc[T] {
	return cursor.NewOffsetAdapter(NewOffsetFinder(db, opts...))
}
