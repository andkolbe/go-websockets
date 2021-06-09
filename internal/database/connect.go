package database

import (
	"github.com/andkolbe/go-websockets/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// create this variable so we can use the db connection inside of the PostRegister handler
var DB *gorm.DB

func Connect(dbConnect string) {
	conn, err := gorm.Open(mysql.Open(dbConnect), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to db")
	}

	DB = conn

	conn.AutoMigrate(&models.User{}, &models.PasswordReset{})
}