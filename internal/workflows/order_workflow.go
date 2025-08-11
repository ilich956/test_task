package workflows

import (
	"fmt"
	"time"

	"forte/internal/models"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflow(ctx workflow.Context, order *models.Order) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("OrderWorkflow started", "OrderID", order.OrderID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	err := workflow.SetQueryHandler(ctx, "GetOrderStatus", func() (*models.Order, error) {
		return order, nil
	})
	if err != nil {
		logger.Error("Failed to set query handler", "Error", err)
		return err
	}

	cancelCh := workflow.GetSignalChannel(ctx, "CancelOrder")

	//Check inventory
	order.UpdateStatus(models.OrderStatusChecking)
	logger.Info("Action 1 Checking inventory", "OrderID", order.OrderID)

	var inventoryResult bool
	err = workflow.ExecuteActivity(ctx, "CheckInventoryActivity", order.Items).Get(ctx, &inventoryResult)
	if err != nil {
		logger.Error("Inventory check failed", "Error", err)
		order.UpdateStatus(models.OrderStatusFailed)
		return fmt.Errorf("inventory check failed: %w", err)
	}

	if !inventoryResult {
		logger.Info("Inventory not available", "OrderID", order.OrderID)
		order.UpdateStatus(models.OrderStatusFailed)
		return fmt.Errorf("insufficient inventory for order %s", order.OrderID)
	}

	//Process payment
	order.UpdateStatus(models.OrderStatusProcessing)
	logger.Info("Action 2 Processing payment", "OrderID", order.OrderID)

	var paymentResult bool
	paymentSelector := workflow.NewSelector(ctx)

	paymentFuture := workflow.ExecuteActivity(ctx, "ProcessPaymentActivity", order.CustomerID, order.TotalAmount)
	paymentSelector.AddFuture(paymentFuture, func(f workflow.Future) {
		err := f.Get(ctx, &paymentResult)
		if err != nil {
			logger.Error("Payment processing failed", "Error", err)
			order.UpdateStatus(models.OrderStatusFailed)
		}
	})

	paymentSelector.AddReceive(cancelCh, func(c workflow.ReceiveChannel, more bool) {
		var reason string
		c.Receive(ctx, &reason)
		logger.Info("Order cancellation requested", "OrderID", order.OrderID, "Reason", reason)
		order.UpdateStatus(models.OrderStatusCancelled)
		paymentResult = false
	})

	paymentSelector.Select(ctx)

	if order.Status == models.OrderStatusCancelled {
		logger.Info("Order cancelled during payment", "OrderID", order.OrderID)
		return nil
	}

	if !paymentResult {
		logger.Info("Payment failed", "OrderID", order.OrderID)
		order.UpdateStatus(models.OrderStatusFailed)
		return fmt.Errorf("payment failed for order %s", order.OrderID)
	}

	//Send notification
	order.UpdateStatus(models.OrderStatusNotifying)
	logger.Info("Action 3 Sending notification", "OrderID", order.OrderID)

	message := fmt.Sprintf("Your order %s processed successfully", order.OrderID)
	err = workflow.ExecuteActivity(ctx, "NotifyCustomerActivity", order.CustomerID, message).Get(ctx, nil)
	if err != nil {
		logger.Error("Notification failed", "Error", err)
	}

	order.UpdateStatus(models.OrderStatusCompleted)
	logger.Info("OrderWorkflow completed successfully", "OrderID", order.OrderID)

	return nil
}
