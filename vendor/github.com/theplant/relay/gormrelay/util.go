package gormrelay

import (
	"reflect"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func parseSchema(db *gorm.DB, v any) (*schema.Schema, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(v); err != nil {
		return nil, errors.Wrap(err, "failed to parse schema for model")
	}
	return stmt.Schema, nil
}

// If T is not a struct or struct pointer, we need to use db.Statement.Model to find or count
func shouldBasedOnModel[T any](db *gorm.DB) (bool, error) {
	if db.Statement.Model != nil {
		return true, nil
	}
	rt := reflect.TypeOf((*T)(nil)).Elem()
	if rt.Kind() == reflect.Struct || (rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct) {
		return false, nil
	}
	return false, errors.New("invalid model type: db.Statement.Model is nil and T is not a struct or struct pointer")
}

func applyModel[T any](db *gorm.DB) *gorm.DB {
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() == reflect.Ptr && reflect.ValueOf(t).IsNil() {
		t = reflect.New(modelType.Elem()).Interface().(T)
	}
	return db.Model(t)
}
