package models

import (
	"github.com/qor5/x/v3/login"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name string

	login.UserPass
	login.SessionSecure
}
