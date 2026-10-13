package models

import "time"

type Expense struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
	CategoryID  int       `json:"category_id"`
}
