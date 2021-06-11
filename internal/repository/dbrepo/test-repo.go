package dbrepo

import (
	"github.com/andkolbe/go-websockets/internal/models"
)

// returns a user by ID
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}