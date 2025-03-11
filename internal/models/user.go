package models

type User struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Balance int64  `json:"balance"`
}
