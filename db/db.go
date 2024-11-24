package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	database, err := gorm.Open(sqlite.Open("bitpin-tui.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	err = database.AutoMigrate(&Session{}, &Favorite{})

	if err != nil {
		panic("failed to migrate database")
	}

	DB = database
}
