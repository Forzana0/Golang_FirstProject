package models

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}
