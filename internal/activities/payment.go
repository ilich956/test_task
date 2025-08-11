package activities

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

type PaymentActivities struct{}

func NewPaymentActivities() *PaymentActivities {
	return &PaymentActivities{}
}

func (a *PaymentActivities) ProcessPaymentActivity(ctx context.Context, customerID string, amount float64) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProcessPaymentActivity started", "CustomerID", customerID, "Amount", amount)

	time.Sleep(3 * time.Second)

	if rand.Float32() < 0.5 {
		logger.Info("Payment failed", "CustomerID", customerID, "Amount", amount)
		return false, fmt.Errorf("payment declined for customer %s", customerID)
	}

	logger.Info("ProcessPaymentActivity completed successfully", "CustomerID", customerID, "Amount", amount)
	return true, nil
}
