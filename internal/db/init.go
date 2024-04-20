package db

import (
	"fmt"

	global "github.com/IUnlimit/telegram-bot-disrecall/internal"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Instance *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s/sqlite.db", global.ParentPath)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database(dsn: %s)", dsn)
	}
	db.AutoMigrate(&model.FileModel{})
	Instance = db
	log.Infof("Initialization of database (dsn: %s) successful", dsn)
}
