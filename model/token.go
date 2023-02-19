package model

import (
	"github.com/google/uuid"
)

type Token struct {
	ID     uuid.UUID `gorm:"primary_key;type:uuid"`
	UserId uuid.UUID `json:"-"`
	Token  string    `json:"-"`
}
