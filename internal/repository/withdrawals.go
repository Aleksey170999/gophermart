package repository

import "github.com/Aleksey170999/go-loyaty/internal/models"

type Withdrawals interface {
	CreateWithdrawal(withdrawal models.Withdrawal) error
	GetWithdrawals(userID int) ([]models.Withdrawal, error)
	GetWithdrawnSum(userID int) (float64, error)
}
