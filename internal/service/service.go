package service

import (
	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Orders interface {
	CreateOrder(number string, userID int) (int, error)
	GetOrdersList(userID int) ([]models.Order, error)
	GetUserBalance(userID int) (int, error)
	GetUserWithdrawn(userID int) (float64, error)
	GetOrderInfo(orderNumber string) (*models.Order, error)
}

type Withdrawals interface {
	ProcessWithdrawal(userID int, orderNumber string, sum int) error
	GetWithdrawals(userID int) ([]models.Withdrawal, error)
	GetWithdrawnSum(userID int) (float64, error)
}

type Service struct {
	Authorization
	Orders
	Withdrawals
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Orders:        NewOrdersService(repos.Orders),
		Withdrawals:   NewWithdrawalsService(repos.Withdrawals),
	}
}
