package models

import "gorm.io/gorm"

type Placement struct {
	gorm.Model

	Name   string
	Active bool

	Units []Unit `gorm:"many2many:placement_units;"`
}
