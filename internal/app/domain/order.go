package domain

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in-progress"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID         int64
	CreatedAt  time.Time
	UserID     string
	TotalPrice int32
	Status     OrderStatus
	Version    int32
}
