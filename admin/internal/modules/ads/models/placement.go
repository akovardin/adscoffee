package models

import "gorm.io/gorm"

type Placement struct {
	gorm.Model

	Name string

	Units []Unit `gorm:"many2many:placement_units;"`
}
