package repository

import (
	"fmt"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

func (r *OrdersPostgres) CreateOrder(order models.Order) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (number, status, accural, updated_at, created_at, user_id) values ($1, $2, $3, $4, $5, $6) RETURNING id", ordersTable)

	row := r.db.QueryRow(query, order.Number, order.Status, order.Accural, order.UpdatedAt, order.CreatedAt, order.UserID)
	if err := row.Scan(&id); err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (r *OrdersPostgres) FindOrderByNumber(orderNumber string) (models.Order, bool, error) {
	var order models.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE number=$1", ordersTable)
	err := r.db.Get(&order, query, orderNumber)
	if err != nil {
		return order, false, err
	}
	return order, true, nil
}

func (r *OrdersPostgres) GetOrdersByUserID(userID int) ([]models.Order, error) {
	var orders []models.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 ORDER BY updated_at", ordersTable)
	err := r.db.Select(&orders, query, userID)
	if err != nil {
		return orders, err
	}
	return orders, nil
}

func (r *OrdersPostgres) GetUserBalance(userID int) (int, error) {
	var balance int
	query := fmt.Sprintf("SELECT current_balance FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&balance, query, userID)
	return balance, err
}

func (r *OrdersPostgres) GetUserWithdrawn(userID int) (float64, error) {
	var withdrawn float64
	query := fmt.Sprintf("SELECT withdrawn_balance FROM %s WHERE user_id=$1 AND accural<0", usersTable)
	err := r.db.Get(&withdrawn, query, userID)
	return withdrawn, err
}

func (r *OrdersPostgres) GetOrderInfo(orderNumber string) (*models.Order, error) {
	var order models.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE number = $1", ordersTable)
	err := r.db.Get(&order, query, orderNumber)
	return &order, err
}
