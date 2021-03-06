package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andkolbe/go-websockets/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// returns a user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	// if a user loses their connection to the internet while in the middle of submitting data to the db, we want that to close and not go through
	// this is called a transaction
	// Go uses something called context to fix this
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // cancel transaction if it takes longer than 3 seconds to complete
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password
		FROM users
		WHERE id = $1
	`

	// QueryRow on its own doesn't know about context. use QueryRowContext instead
	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
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

// updates user in the database
func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second) // cancel transaction if it takes longer than 3 seconds to complete
	defer cancel()

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3
	`
	_, err := m.DB.ExecContext(ctx, query,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)
	if err != nil {
		return err
	}

	return nil
}

// Register
func (m *postgresDBRepo) Register(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second) // cancel transaction if it takes longer than 3 seconds to complete
	defer cancel()

	// create a bcrypt hash of the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, 12)
	if err != nil {
		return 0, err
	}

	var newId int

	query := `
		INSERT INTO users (first_name, last_name, email, password) 
		VALUES ($1, $2, $3, $4) returning id
	`

	err = m.DB.QueryRowContext(ctx, query, 
		&user.FirstName,
		&user.LastName,
		&user.Email,
		hashedPassword,
	).Scan(&newId)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return newId, err
}

// Authenticate
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // cancel transaction if it takes longer than 3 seconds to complete
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

// func (m *mySQLDBRepo) ResetPassword(id int, newPassword string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Create a bcrypt hash of the plain-text password.
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	query := `UPDATE users SET password = $1 WHERE id = $2`

// 	_, err = m.DB.ExecContext(ctx, query, hashedPassword, id)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	return nil
// }