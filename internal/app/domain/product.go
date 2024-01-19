package domain

import "time"

type Product struct {
	ID          int64
	CreatedAt   time.Time
	Title       string
	Description string
	Price       int32
	Version     int32
}
