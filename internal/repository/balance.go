package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type BalanceRepository interface {
	UpdateUserBalance(userID int, amount int) error
}

type balanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) BalanceRepository {
	return &balanceRepository{db: db}
}

func (r *balanceRepository) UpdateUserBalance(userID int, amount int) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Update the balance
	_, err = tx.Exec(
		`UPDATE users 
		set current_balance = current_balance + $1 
		WHERE id = $2`,
		amount,
		userID,
	)
	if err != nil {
		return err
	}

	// If this is a withdrawal, also update the spent amount
	if amount < 0 {
		_, err = tx.Exec(
			`UPDATE users 
			set withdrawn_balance = withdrawn_balance + $1 
			WHERE id = $2`,
			-amount, // Convert to positive for spent amount
			userID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
