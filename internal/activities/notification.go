package activities

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type NotificationActivities struct{}

func NewNotificationActivities() *NotificationActivities {
	return &NotificationActivities{}
}

func (a *NotificationActivities) NotifyCustomerActivity(ctx context.Context, customerID string, message string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("NotifyCustomerActivity started", "CustomerID", customerID)

	time.Sleep(3 * time.Second)

	logger.Info("Notification sent", "CustomerID", customerID, "Message", message)

	logger.Info("NotifyCustomerActivity completed successfully")
	return nil
}
