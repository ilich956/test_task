package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"forte/internal/constants"
	"forte/internal/models"

	"go.temporal.io/sdk/client"
)

type CreateOrderRequest struct {
	CustomerID string             `json:"customerId"`
	Items      []models.OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderID    string `json:"orderId"`
	WorkflowID string `json:"workflowId"`
	Status     string `json:"status"`
}

type OrderStatusResponse struct {
	OrderID    string             `json:"orderId"`
	WorkflowID string             `json:"workflowId"`
	Status     string             `json:"status"`
	Items      []models.OrderItem `json:"items,omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty"`
}

type OrderHandlers struct {
	temporalClient client.Client
}

func NewOrderHandlers(temporalClient client.Client) *OrderHandlers {
	return &OrderHandlers{
		temporalClient: temporalClient,
	}
}

func (h *OrderHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.CustomerID == "" || len(req.Items) == 0 {
		http.Error(w, "CustomerID and Items are required", http.StatusBadRequest)
		return
	}

	order := models.NewOrder(req.CustomerID, req.Items)

	workflowID := fmt.Sprintf("order-workflow-%s", order.OrderID)
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: constants.ORDER_TASK_QUEUE,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workflowRun, err := h.temporalClient.ExecuteWorkflow(ctx, workflowOptions, constants.ORDER_WORKFLOW, order)
	if err != nil {
		log.Printf("Failed to start workflow: %v", err)
		http.Error(w, "Failed to start order processing", http.StatusInternalServerError)
		return
	}

	response := CreateOrderResponse{
		OrderID:    order.OrderID,
		WorkflowID: workflowRun.GetID(),
		Status:     string(models.OrderStatusCreated),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Printf("Started order workflow: OrderID=%s, WorkflowID=%s", order.OrderID, workflowRun.GetID())
}

func (h *OrderHandlers) HandleOrder(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	orderID := parts[0]

	if len(parts) == 1 && r.Method == http.MethodGet {
		h.getOrderStatus(w, orderID)
	} else if len(parts) == 2 && parts[1] == "cancel" && r.Method == http.MethodPost {
		h.cancelOrder(w, orderID)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *OrderHandlers) getOrderStatus(w http.ResponseWriter, orderID string) {
	workflowID := fmt.Sprintf("order-workflow-%s", orderID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.temporalClient.QueryWorkflow(ctx, workflowID, "", "GetOrderStatus")
	if err != nil {
		log.Printf("Failed to query workflow %s: %v", workflowID, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	var order models.Order
	if err := resp.Get(&order); err != nil {
		log.Printf("Failed to decode order status: %v", err)
		http.Error(w, "Failed to get order status", http.StatusInternalServerError)
		return
	}

	response := OrderStatusResponse{
		OrderID:    order.OrderID,
		WorkflowID: workflowID,
		Status:     string(order.Status),
		Items:      order.Items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandlers) cancelOrder(w http.ResponseWriter, orderID string) {
	workflowID := fmt.Sprintf("order-workflow-%s", orderID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.temporalClient.SignalWorkflow(ctx, workflowID, "", "CancelOrder", "User requested cancellation")
	if err != nil {
		log.Printf("Failed to send cancel signal to workflow %s: %v", workflowID, err)
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Cancellation request sent",
		"orderId": orderID,
	})

	log.Printf("Sent cancellation signal to order: %s", orderID)
}
