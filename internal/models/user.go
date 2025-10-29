package models

type User struct {
	ID        int    `json:"id" db:"id"`
	UserName  string `json:"username" db:"username"`
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Password  string `json:"password" db:"password"`
	Balance   int    `json:"current" db:"current_balance"`
	Spent     int    `json:"withdrawn,omitempty" db:"withdrawn_balance"`
}
