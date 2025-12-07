package entity

import "time"

type Topic struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Count     int       `db:"count"`
	CreatedAt time.Time `db:"created_at"`
}
