package service

import (
	"time"

	"github.com/Aleksey170999/go-loyaty/internal/apperror"
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

// Order status constants
const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

func (s *OrdersService) CreateOrder(orderNumber string, userID int) (int, error) {

	order, found, err := s.repo.FindOrderByNumber(orderNumber)
	if err != nil {
		return 0, apperror.Wrap(err, "failed to find order", 0)
	}

	if found {
		if order.UserID == userID {
			return 0, apperror.ErrOrderByUserExists
		}
		return 0, apperror.ErrOrderByOtherExists
	}

	newOrder := models.Order{
		Number:    orderNumber,
		Status:    StatusNew,
		Accural:   0,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		UserID:    userID,
	}

	id, err := s.repo.CreateOrder(newOrder)
	if err != nil {
		return 0, apperror.Wrap(err, "failed to create order", 0)
	}

	return id, nil
}

func (s *OrdersService) GetOrdersList(userID int) ([]models.Order, error) {
	orders, err := s.repo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, apperror.Wrap(err, "failed to get orders list", 0)
	}
	return orders, nil
}

func (s *OrdersService) GetOrderInfo(orderNumber string) (*models.Order, error) {
	order, err := s.repo.GetOrderInfo(orderNumber)
	if err != nil {
		return nil, apperror.Wrap(err, "failed to get order info", 0)
	}
	return order, nil
}
