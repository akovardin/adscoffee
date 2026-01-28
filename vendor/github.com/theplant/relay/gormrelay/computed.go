package gormrelay

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/theplant/relay/cursor"
)

// Computed defines SQL expressions calculated at database level and attached to query results.
type Computed[T any] struct {
	// Maps field names to SQL expressions
	Columns map[string]clause.Column

	// Function for scanning rows and transforming results
	ForScan func(db *gorm.DB) (dest any, toCursorNodes func(computedResults []map[string]any) []cursor.Node[T], err error)
}

// Validate ensures the Computed configuration is valid.
func (c *Computed[T]) Validate() error {
	if len(c.Columns) == 0 {
		return errors.New("Columns must not be empty")
	}

	for field, col := range c.Columns {
		if col.Alias != "" {
			return errors.Errorf("Column %q should have empty Alias, got %q", field, col.Alias)
		}
	}

	if c.ForScan == nil {
		return errors.New("ForScan function must not be nil")
	}

	aliasMap := make(map[string][]string)
	for fieldName := range c.Columns {
		alias := ComputedFieldToColumnAlias(fieldName)
		aliasMap[alias] = append(aliasMap[alias], fieldName)
	}

	var duplicates []string
	for alias, fields := range aliasMap {
		if len(fields) > 1 {
			duplicates = append(duplicates, fmt.Sprintf("alias %q is used by multiple fields: %v", alias, fields))
		}
	}

	if len(duplicates) > 0 {
		return errors.Errorf("duplicate computed field aliases detected:\n%s", strings.Join(duplicates, "\n"))
	}

	return nil
}

// WithComputed adds computed fields to a query.
func WithComputed[T any](computed *Computed[T]) Option[T] {
	return func(opts *Options[T]) {
		opts.Computed = computed
	}
}

// ComputedColumns creates a map of Column objects from field-to-SQL expression mappings.
var ComputedColumns = func(columns map[string]string) map[string]clause.Column {
	return lo.MapEntries(columns, func(field string, value string) (string, clause.Column) {
		value = strings.Trim(value, " ()")
		return field, clause.Column{Name: fmt.Sprintf("(%s)", value), Raw: true}
	})
}

// ComputedFieldToColumnAlias generates a standardized SQL alias for computed fields.
var ComputedFieldToColumnAlias = func(field string) string {
	return strings.ToLower("_relay_computed_" + lo.SnakeCase(field))
}

// withComputedResult combines an object with computed field values for JSON serialization.
type withComputedResult struct {
	Object          any
	ComputedResults map[string]any
}

// MarshalJSON implements json.Marshaler interface.
func (v *withComputedResult) MarshalJSON() ([]byte, error) {
	b, err := cursor.JSONMarshal(v.Object)
	if err != nil {
		return nil, err
	}

	for path, value := range v.ComputedResults {
		b, err = sjson.SetBytes(b, path, value)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

// WithComputedResult creates a wrapper that merges an object with computed field values.
func WithComputedResult(object any, computedResults map[string]any) *withComputedResult {
	return &withComputedResult{
		Object:          object,
		ComputedResults: computedResults,
	}
}

// DefaultForScan creates a standard ForScan function for Computed.
// This simplifies the creation of a ForScan function that wraps objects with
// their computed values and makes them accessible in the resulting cursor.Node objects.
func DefaultForScan[T any](db *gorm.DB) (dest any, toCursorNodes func(computedResults []map[string]any) []cursor.Node[T], err error) {
	if db.Statement.Model == nil {
		return nil, nil, errors.New("db.Statement.Model cannot be nil")
	}

	modelType := reflect.TypeOf(db.Statement.Model)
	genericType := reflect.TypeOf((*T)(nil)).Elem()

	// Only use []T when Model type exactly matches T
	if modelType == genericType {
		nodes := []T{}
		return &nodes, func(computedResults []map[string]any) []cursor.Node[T] {
			return lo.Map(nodes, func(node T, i int) cursor.Node[T] {
				return &cursor.NodeWrapper[T]{
					Object: WithComputedResult(node, computedResults[i]),
					Unwrap: func() T { return node },
				}
			})
		}, nil
	}

	sliceType := reflect.SliceOf(modelType)
	nodesVal := reflect.New(sliceType).Elem()

	return nodesVal.Addr().Interface(), func(computedResults []map[string]any) []cursor.Node[T] {
		result := make([]cursor.Node[T], nodesVal.Len())
		for i := 0; i < nodesVal.Len(); i++ {
			node := nodesVal.Index(i).Interface().(T) // Direct type assertion - panic is expected if types are incompatible
			result[i] = &cursor.NodeWrapper[T]{
				Object: WithComputedResult(node, computedResults[i]),
				Unwrap: func() T { return node },
			}
		}
		return result
	}, nil
}

func splitComputedScan(computedColumns map[string]clause.Column, tx *gorm.DB, dest any) (computedResults []map[string]any, err error) {
	aliasToField := make(map[string]string)
	splitter := make(map[string]func(columnType *sql.ColumnType) any)
	for field := range computedColumns {
		alias := ComputedFieldToColumnAlias(field)
		aliasToField[alias] = field
		splitter[alias] = func(columnType *sql.ColumnType) any {
			return columnType.ScanType()
		}
	}

	if err := SplitScan(tx, dest, splitter, &computedResults).Error; err != nil {
		return nil, errors.Wrap(err, "failed to scan records with computed columns")
	}

	for i, result := range computedResults {
		computedResults[i] = lo.MapEntries(result, func(alias string, value any) (string, any) {
			return aliasToField[alias], value
		})
	}

	return computedResults, nil
}
