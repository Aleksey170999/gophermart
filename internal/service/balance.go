package service

import (
	"errors"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
)

type BalanceService struct {
	repo repository.BalanceRepository
}

func NewBalanceService(repo repository.BalanceRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

// UpdateUserBalance updates user's balance
func (s *BalanceService) UpdateUserBalance(userID int, amount int) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	return s.repo.UpdateUserBalance(userID, amount)
}

// GetUserBalance retrieves the current balance for a user
func (s *BalanceService) GetUserBalance(userID int) (*models.User, error) {
	// This method should be implemented in the repository layer
	// For now, we'll return an error indicating it's not implemented
	return nil, errors.New("not implemented: use repository layer to get user balance")
}
