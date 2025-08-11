package models

import (
	"time"

	"github.com/google/uuid"
)

type orderStatus string

const (
	OrderStatusCreated    orderStatus = "created"
	OrderStatusChecking   orderStatus = "checking"
	OrderStatusProcessing orderStatus = "processing"
	OrderStatusNotifying  orderStatus = "notifying"

	OrderStatusCompleted orderStatus = "completed"
	OrderStatusFailed    orderStatus = "failed"
	OrderStatusCancelled orderStatus = "cancelled"
)

type Order struct {
	OrderID     string
	CustomerID  string
	Items       []OrderItem
	TotalAmount float64

	Status    orderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ProductID string
	Name      string
	Price     float64
	Quantity  int
}

func NewOrder(customerID string, items []OrderItem) *Order {
	order := &Order{
		OrderID:    uuid.NewString(),
		CustomerID: customerID,
		Items:      items,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	var total float64
	for _, i := range items {
		total = total + i.Price*float64(i.Quantity)
	}

	order.TotalAmount = total
	order.Status = OrderStatusCreated

	return order
}

func (o *Order) UpdateStatus(status orderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}
