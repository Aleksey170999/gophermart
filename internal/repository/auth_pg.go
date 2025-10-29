package repository

import (
	"database/sql"
	"fmt"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, first_name, last_name, password) values ($1, $2, $3, $4) RETURNING id", usersTable)

	err := r.db.QueryRow(query, user.UserName, user.FirstName, user.LastName, user.Password).Scan(&id)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return 0, fmt.Errorf("user with this username already exists")
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, username, password, first_name, last_name FROM %s WHERE username = $1", usersTable)
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUser is kept for backward compatibility but should not be used in new code
// Instead, use GetUserByUsername and compare password hashes in the service layer
func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	// This method is now just a wrapper that doesn't verify the password
	// The actual password check should be done in the service layer using bcrypt
	return user, nil
}
