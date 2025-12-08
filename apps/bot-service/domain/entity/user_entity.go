package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `db:"id"`
	PhoneNumber string     `db:"phone_number"`
	Name        string     `db:"name"`
	JobTitle    *string    `db:"job_title"`
	Gender      *string    `db:"gender"`
	DateOfBirth *time.Time `db:"date_of_birth"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type GetUsersFilter struct {
	Offset int
	Limit  int
	Search string
}
