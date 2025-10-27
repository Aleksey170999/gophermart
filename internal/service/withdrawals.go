package service

import (
	"fmt"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
)

type WithdrawalsService struct {
	repo repository.Withdrawals
}

func NewWithdrawalsService(repo repository.Withdrawals) *WithdrawalsService {
	return &WithdrawalsService{repo: repo}
}

func (s *WithdrawalsService) ProcessWithdrawal(userID int, orderNumber string, sum int) error {
	if sum <= 0 {
		return fmt.Errorf("sum must be positive")
	}

	withdrawal := models.Withdrawal{
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
	}

	err := s.repo.CreateWithdrawal(withdrawal)
	if err != nil {
		return fmt.Errorf("error processing withdrawal: %v", err)
	}

	return nil
}

func (s *WithdrawalsService) GetWithdrawals(userID int) ([]models.Withdrawal, error) {
	return s.repo.GetWithdrawals(userID)
}

func (s *WithdrawalsService) GetWithdrawnSum(userID int) (float64, error) {
	return s.repo.GetWithdrawnSum(userID)
}
