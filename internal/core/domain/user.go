package domain

import (
	"github.com/gostructure/app/pkg/util/time_util"
)

type User struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"-"` // Never return password hash in JSON
	Role         string             `json:"role"`
	Status       string             `json:"status"`
	TokenVersion int                `json:"-"`
	CreatedAt    time_util.JSONTime `json:"created_at"`
	UpdatedAt    time_util.JSONTime `json:"updated_at"`
}
