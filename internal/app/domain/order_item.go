package domain

type OrderItem struct {
	OrderID   int64
	ProductID int64
	Quantity  int32
	Price     int32
	Version   int32
}
