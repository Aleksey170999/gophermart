package repository

import (
	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}
type Orders interface {
	CreateOrder(models.Order) (int, error)
	FindOrderByNumber(orderNumber string) (models.Order, bool, error)
	GetOrdersByUserID(userID int) ([]models.Order, error)
	GetUserBalance(userID int) (int, error)
	GetUserWithdrawn(userID int) (float64, error)
	GetOrderInfo(orderNumber string) (*models.Order, error)
}

type Repository struct {
	Authorization
	Orders
	Withdrawals
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Orders:        NewOrdersPostgres(db),
		Withdrawals:   NewWithdrawalsPostgres(db),
	}
}
