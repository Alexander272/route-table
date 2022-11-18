package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Login    string    `db:"login" json:"login"`
	Password string    `db:"password" json:"-"`
	Role     string    `db:"role" json:"role"`
	// Operations pq.StringArray `db:"Operations" json:"operations"`
}

type UserDTO struct {
	Id       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	RoleId   uuid.UUID `json:"roleId"`
}
