package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/jmoiron/sqlx"
)

type BalanceService struct {
	db *sqlx.DB
}

func NewBalanceService(db *sqlx.DB) *BalanceService {
	return &BalanceService{db: db}
}

// UpdateUserBalance updates user's balance atomically
func (s *BalanceService) UpdateUserBalance(userID int, amount int) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Update the balance
	_, err = tx.Exec(
		`UPDATE users 
		SET current_balance = current_balance + $1 
		WHERE id = $2`,
		amount,
		userID,
	)
	if err != nil {
		return fmt.Errorf("error updating balance: %v", err)
	}

	// If this is a withdrawal, also update the spent amount
	if amount < 0 {
		_, err = tx.Exec(
			`UPDATE users 
			SET withdrawn_balance = withdrawn_balance + $1 
			WHERE id = $2`,
			-amount, // Convert to positive for spent amount
			userID,
		)
		if err != nil {
			return fmt.Errorf("error updating spent amount: %v", err)
		}
	}

	return tx.Commit()
}

// GetUserBalance retrieves the current balance for a user
func (s *BalanceService) GetUserBalance(userID int) (*models.User, error) {
	var user models.User
	err := s.db.Get(&user, 
		`SELECT id, username, current_balance, withdrawn_balance 
		 FROM users WHERE id = $1`, 
		userID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user balance: %v", err)
	}

	return &user, nil
}
