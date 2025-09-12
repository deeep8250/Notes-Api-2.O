package db

import (
	"fmt"
	"pr01/config"
	"pr01/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbInit(cfg config.Config) {

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		fmt.Println("db error : ", err.Error())
		return
	}

	db.AutoMigrate(models.User{}, models.Notes{})
	DB = db

}
