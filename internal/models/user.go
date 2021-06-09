package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName  string `json:"username"` // format the JSON response to look like these values instead
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  []byte `json:"-"`
}