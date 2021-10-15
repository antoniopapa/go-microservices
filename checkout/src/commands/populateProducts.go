package main

import (
	"checkout/src/database"
	"checkout/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	database.Connect()

	db, err := gorm.Open(mysql.Open("root:root@tcp(host.docker.internal:33066)/ambassador"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	var products []models.Product

	db.Find(&products)

	for _, product := range products {
		database.DB.Create(&product)
	}
}
