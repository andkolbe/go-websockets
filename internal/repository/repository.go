package repository

import "github.com/andkolbe/go-websockets/internal/models"

type DatabaseRepo interface {
	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Register(user models.User) (int, error)
	Login(username, testPassword string) (int, string, error)
}