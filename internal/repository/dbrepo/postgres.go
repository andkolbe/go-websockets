package dbrepo

import (
	"context"
	"time"

	"github.com/andkolbe/go-websockets/internal/models"
)

// returns a user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	// if a user loses their connection to the internet while in the middle of submitting data to the db, we want that to close and not go through
	// this is called a transaction
	// Go uses something called context to fix this
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // cancel transaction if it takes longer than 3 seconds to complete
	defer cancel()

	query := `
		SELECT id, username, first_name, last_name, email, password
		FROM users
		WHERE id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second) // cancel transaction if it takes longer than 3 seconds to complete
	defer cancel()

	query := `
		UPDATE users
		SET username = $1, first_name = $2, last_name = $3, email = $4
	`
	_, err := m.DB.ExecContext(ctx, query,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)
	if err != nil {
		return err
	}

	return nil
}