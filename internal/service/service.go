package service

import (
	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type Authorization interface {
	// Core methods for the repository
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	
	// New methods for the service layer
	Login(username, password string) (AuthResponse, error)
	Register(user models.User) (AuthResponse, error)
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
	JWTService *JWTService
}

func NewService(repos *repository.Repository, jwtSecret string) *Service {
	jwtService := NewJWTService(jwtSecret)
	
	return &Service{
		Authorization: NewAuthService(repos.Authorization, jwtService),
		Orders:        NewOrdersService(repos.Orders),
		Withdrawals:   NewWithdrawalsService(repos.Withdrawals),
		JWTService:    jwtService,
	}
}
