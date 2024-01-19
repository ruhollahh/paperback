package domain

import "time"

type InvoiceStatus string

const (
	InvoiceStatusUnpaid InvoiceStatus = "unpaid"
	InvoiceStatusPaid   InvoiceStatus = "paid"
)

type Invoice struct {
	ID        int64
	OrderID   int64
	CreatedAt time.Time
	Status    InvoiceStatus
	Version   int32
}
