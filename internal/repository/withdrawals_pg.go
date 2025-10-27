package repository

import (
	"fmt"
	"time"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/jmoiron/sqlx"
)

const withdrawalsTable = "withdrawals"

type WithdrawalsPostgres struct {
	db *sqlx.DB
}

func NewWithdrawalsPostgres(db *sqlx.DB) *WithdrawalsPostgres {
	return &WithdrawalsPostgres{db: db}
}

var _ Withdrawals = (*WithdrawalsPostgres)(nil)

func (r *WithdrawalsPostgres) CreateWithdrawal(withdrawal models.Withdrawal) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if the order already exists
	var exists bool
	err = tx.Get(&exists,
		"SELECT EXISTS(SELECT 1 FROM withdrawals WHERE order_number = $1)",
		withdrawal.OrderNumber,
	)
	if err != nil {
		return fmt.Errorf("error checking order existence: %v", err)
	}
	if exists {
		return fmt.Errorf("order already exists")
	}

	query := `select current_balance from users where id = $1`

	var availableBalance int
	err = tx.Get(&availableBalance, query, withdrawal.UserID)
	if err != nil {
		return fmt.Errorf("error getting available balance: %v", err)
	}

	// Check if there's enough balance
	if availableBalance < withdrawal.Sum {
		return fmt.Errorf("insufficient funds: available %d, requested %d", availableBalance, withdrawal.Sum)
	}

	_, err = tx.Exec(
		`UPDATE users 
		SET current_balance = current_balance - $1 
		WHERE id = $2`,
		withdrawal.Sum,
		withdrawal.UserID,
	)
	if err != nil {
		return fmt.Errorf("error updating user balance: %v", err)
	}

	// Create the withdrawal record
	_, err = tx.Exec(
		`INSERT INTO withdrawals (user_id, order_number, sum, processed_at) 
		 VALUES ($1, $2, $3, $4)`,
		withdrawal.UserID,
		withdrawal.OrderNumber,
		withdrawal.Sum,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("error creating withdrawal: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (r *WithdrawalsPostgres) GetWithdrawals(userID int) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	query := fmt.Sprintf(`SELECT id, user_id, order_number, sum, processed_at FROM %s WHERE user_id = $1 ORDER BY processed_at DESC`, withdrawalsTable)

	err := r.db.Select(&withdrawals, query, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (r *WithdrawalsPostgres) GetWithdrawnSum(userID int) (float64, error) {
	var sum float64
	query := fmt.Sprintf(`SELECT COALESCE(SUM(sum), 0) FROM %s WHERE user_id = $1`, withdrawalsTable)

	err := r.db.Get(&sum, query, userID)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("error getting withdrawn sum: %v", err)
	}

	return sum, nil
}
