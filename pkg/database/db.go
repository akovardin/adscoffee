package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(config Config) *gorm.DB {
	connection := config.Connection()

	// Create database connection
	db, err := gorm.Open(postgres.Open(connection))
	if err != nil {
		panic(err)
	}

	// Set db log level
	if config.Debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Warn)
	}

	return db
}
