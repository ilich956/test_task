package activities

type OrderActivities struct {
	Inventory    *InventoryActivities
	Payment      *PaymentActivities
	Notification *NotificationActivities
}

func NewOrderActivities() *OrderActivities {
	return &OrderActivities{
		Inventory:    NewInventoryActivities(),
		Payment:      NewPaymentActivities(),
		Notification: NewNotificationActivities(),
	}
}

func (a *OrderActivities) GetActivities() []interface{} {
	return []interface{}{
		a.Inventory.CheckInventoryActivity,
		a.Payment.ProcessPaymentActivity,
		a.Notification.NotifyCustomerActivity,
	}
}
