package db

import (
	"bitpin-tui/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func initiateSettings() {
	InsertIfNotExistsSetting("base_url", "https://api.bitpin.market")
	UpsertSetting("version", utils.GetVersion())
	InsertIfNotExistsSetting("pussy_out_workers", "10")
	InsertIfNotExistsSetting("pussy_out_buffer_size", "100")
}

func init() {
	database, err := gorm.Open(sqlite.Open("bitpin-tui.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	err = database.AutoMigrate(&Session{}, &Favorite{}, &Setting{})

	if err != nil {
		panic("failed to migrate database")
	}

	DB = database

	initiateSettings()
}
