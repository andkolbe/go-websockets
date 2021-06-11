package models

import "time"

// we have to describe our database in a way that Go understands
// every table gets a model
type User struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  []byte
	CreatedAt time.Time
}

/*
The primitive data types prefixed with "u" are unsigned versions with the same bit sizes.
Effectively, this means they cannot store negative numbers, but on the other hand they can store positive numbers twice as large as their signed counterparts.
The signed counterparts do not have "u" prefixed.
*/
