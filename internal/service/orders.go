package service

import (
	"errors"
	"time"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/Aleksey170999/go-loyaty/internal/repository"
)

type OrdersService struct {
	repo repository.Orders
}

func (s *OrdersService) GetUserBalance(userID int) (int, error) {
	return s.repo.GetUserBalance(userID)
}

func (s *OrdersService) GetUserWithdrawn(userID int) (float64, error) {
	return s.repo.GetUserWithdrawn(userID)
}

func NewOrdersService(repo repository.Orders) *OrdersService {
	return &OrdersService{repo: repo}

}

type OrderStatus int

const (
	New        string = "NEW"
	Processing string = "PROCESSING"
	Invalid    string = "INVALID"
	Processed  string = "PROCESSED"
)

var (
	ErrInvalidOrderNumber = errors.New("invalid order number format")
	ErrOrderByUserExists  = errors.New("order already uploaded by this user")
	ErrOrderByOtherExists = errors.New("order already uploaded by another user")
)

func (s *OrdersService) CreateOrder(orderNumber string, userID int) (int, error) {
	order, found, err := s.repo.FindOrderByNumber(orderNumber)
	if err != nil && found {
		return 0, err
	}
	if found {
		if order.UserID == userID {
			return 0, ErrOrderByUserExists
		} else {
			return 0, ErrOrderByOtherExists
		}
	}
	newOrder := models.Order{
		Number:    orderNumber,
		Status:    New,
		Accural:   0,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		UserID:    userID,
	}
	return s.repo.CreateOrder(newOrder)
}

func (s *OrdersService) GetOrdersList(userID int) ([]models.Order, error) {
	orders, err := s.repo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrdersService) GetOrderInfo(orderNumber string) (*models.Order, error) {
	order, err := s.repo.GetOrderInfo(orderNumber)
	if err != nil {
		return nil, err
	}
	return order, nil
}
