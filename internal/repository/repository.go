package repository

import "github.com/andkolbe/go-websockets/internal/models"

type DatabaseRepo interface {
	GetUserByID(id int) (models.User, error)
}