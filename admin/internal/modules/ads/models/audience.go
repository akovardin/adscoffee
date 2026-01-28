package models

import "gorm.io/gorm"

const (
	StatusPending   = "pending"
	StatusProcessed = "processed"
)

type Audience struct {
	gorm.Model

	Title string
	Name  string

	File   string
	Status string
}
