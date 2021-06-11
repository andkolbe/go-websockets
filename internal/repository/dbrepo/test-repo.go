package dbrepo

import (
	"github.com/andkolbe/go-websockets/internal/models"
)

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

func (m *testDBRepo) Register(user models.User) error {
	return nil
}

func (m *testDBRepo) Login(username, testPassword string) (int, string, error) {
	var id int
	var hashedPassword string
	return id, hashedPassword, nil
}