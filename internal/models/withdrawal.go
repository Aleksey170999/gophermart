package models

import "time"

type WithdrawalRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type Withdrawal struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	OrderNumber string    `db:"order_number"`
	Sum         int       `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}
