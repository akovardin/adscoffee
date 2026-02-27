package models

import "gorm.io/gorm"

type Unit struct {
	gorm.Model

	Name  string
	Price int

	NetworkID int
	Network   Network

	PlacementID int
	Placement   Placement

	Data   string
	Active bool
}
