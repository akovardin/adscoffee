package models

import "gorm.io/gorm"

type Network struct {
	gorm.Model

	Title string
	Name  string
}
