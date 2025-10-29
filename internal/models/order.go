package models

import "time"

type Order struct {
	ID        int       `db:"id" json:"id"`
	Number    string    `db:"number" json:"number"`
	Status    string    `db:"status" json:"status"`
	Accural   int       `db:"accural" json:"accural"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserID    int       `db:"user_id" json:"user_id"`
}
