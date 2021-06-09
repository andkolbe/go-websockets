package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dbConnect string) {
	conn, err := gorm.Open(mysql.Open(dbConnect), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to db")
	}

	DB = conn

	fmt.Println(DB)
}