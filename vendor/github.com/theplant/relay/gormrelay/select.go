package gormrelay

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var SELECT = clause.Select{}.Name()

func AppendSelect(columns ...clause.Column) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(columns) == 0 {
			return db
		}
		clauseSelect, ok := db.Statement.Clauses[SELECT]
		if !ok {
			clauseSelect = clause.Clause{
				Name: SELECT,
			}
		}
		oldBuilder := clauseSelect.Builder
		clauseSelect.Builder = func(c clause.Clause, builder clause.Builder) {
			if exprSelect, ok := c.Expression.(clause.Select); ok {
				if len(exprSelect.Columns) == 0 {
					exprSelect.Columns = make([]clause.Column, 0, len(columns)+1)
					exprSelect.Columns = append(exprSelect.Columns, clause.Column{
						Name: "*",
						Raw:  true,
					})
				}
				exprSelect.Columns = append(exprSelect.Columns, columns...)
				c.Expression = exprSelect
			}

			if oldBuilder != nil {
				oldBuilder(c, builder)
			} else {
				c.Builder = nil
				c.Build(builder)
			}
		}
		db.Statement.Clauses[SELECT] = clauseSelect
		return db
	}
}
