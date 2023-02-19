package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primary_key;type:uuid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}
