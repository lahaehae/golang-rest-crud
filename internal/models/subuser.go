package models

import (
	"time"
)

type SubUser struct {
    ID          string    `json:"id" db:"id"`
    OwnerID     string    `json:"owner_id" db:"owner_id"`     // ID пользователя из Kratos
    Name        string    `json:"name" db:"name"`
    Email       string    `json:"email" db:"email"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}