package activities

import (
	"context"
	"math/rand"
	"time"

	"forte/internal/models"

	"go.temporal.io/sdk/activity"
)

type InventoryActivities struct{}

func NewInventoryActivities() *InventoryActivities {
	return &InventoryActivities{}
}

func (a *InventoryActivities) CheckInventoryActivity(ctx context.Context, items []models.OrderItem) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CheckInventoryActivity started")

	time.Sleep(3 * time.Second)

	for _, item := range items {
		logger.Info("Checking inventory for:", "ProductID", item.ProductID, "Quantity", item.Quantity)
		//change probability
		if rand.Float32() < 0.1 {
			logger.Info("Insufficient inventory", "ProductID", item.ProductID)
			return false, nil
		}
	}

	logger.Info("CheckInventoryActivity completed successfully")
	return true, nil
}
