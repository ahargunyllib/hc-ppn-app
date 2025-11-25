package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id"`
	PhoneNumber string    `db:"phone_number"`
	Label       string    `db:"label"`
	AssignedTo  *string   `db:"assigned_to"`
	Notes       *string   `db:"notes"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type GetUsersFilter struct {
	Offset     int
	Limit      int
	Search     string
	AssignedTo *string
}
