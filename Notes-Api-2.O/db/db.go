package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB gorm.DB

func DbInit() {

	dsn := "host=localhost user=postgres password=deep dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("db error : ", err.Error())
		return
	}

	db.AutoMigrate()
	DB = *db

}
