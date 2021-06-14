package dbrepo

import (
	"errors"

	"github.com/andkolbe/go-websockets/internal/models"
)

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

func (m *testDBRepo) Register(user models.User) (int, error) {
	var newId int
	return newId, nil
}

func (m *testDBRepo) Login(username, testPassword string) (int, string, error) {
	if username == "test" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")
}