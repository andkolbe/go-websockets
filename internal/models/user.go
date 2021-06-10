package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName  string `json:"username"` // format the JSON response to look like these values instead
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  []byte `json:"-"` // we don't want to send the password to the client
}

/*
The primitive data types prefixed with "u" are unsigned versions with the same bit sizes. 
Effectively, this means they cannot store negative numbers, but on the other hand they can store positive numbers twice as large as their signed counterparts. 
The signed counterparts do not have "u" prefixed.
*/
