package models

import (
	"time"

	"gorm.io/gorm"
)

type Advertiser struct {
	gorm.Model

	Title  string
	Info   string
	Active bool

	Start time.Time
	End   time.Time

	Targeting string
	Budget    string
	Capping   string
	Timetable string

	OrdContract string
	OrdEnable   bool
}
