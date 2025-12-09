package entity

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Rating    int       `db:"rating"`
	Comment   *string   `db:"comment"`
	CreatedAt time.Time `db:"created_at"`

	User User `db:"user"`
}

type GetFeedbacksFilter struct {
	Offset    int
	Limit     int
	UserID    *uuid.UUID
	Ratings   []int
	MinRating *int
	MaxRating *int
}

type SatisfactionTrendRow struct {
	Date            time.Time `db:"date"`
	AvgSatisfaction float64   `db:"avg_satisfaction"`
}
