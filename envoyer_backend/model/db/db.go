package db

import (
	"envoyer/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDb() *gorm.DB {

	var level logger.LogLevel
	if config.Config.GinMode == "debug" {
		level = logger.Info
	} else {
		level = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(config.Config.DbUrl), &gorm.Config{PrepareStmt: true, Logger: logger.Default.LogMode(level)})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(config.Config.MaxDbConnection)
	sqlDb.SetMaxOpenConns(config.Config.MaxDbConnection)
	fmt.Println("mysql connection successful...")

	return db

	return &gorm.DB{}
}
