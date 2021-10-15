package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"users/src/database"
	"users/src/models"
)

func main() {
	database.Connect()

	db, err := gorm.Open(mysql.Open("root:root@tcp(host.docker.internal:33066)/ambassador"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	var users []models.User

	db.Find(&users)

	for _, user := range users {
		database.DB.Create(&user)
	}
}
