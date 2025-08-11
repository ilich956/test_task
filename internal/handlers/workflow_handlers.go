package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"go.temporal.io/sdk/client"
)

type WorkflowHandlers struct {
	temporalClient client.Client
}

func NewWorkflowHandlers(temporalClient client.Client) *WorkflowHandlers {
	return &WorkflowHandlers{
		temporalClient: temporalClient,
	}
}

func (h *WorkflowHandlers) GetWorkflowStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workflowID := strings.TrimPrefix(r.URL.Path, "/api/workflows/")
	workflowID = strings.TrimSuffix(workflowID, "/status")

	if workflowID == "" {
		http.Error(w, "Workflow ID required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.temporalClient.DescribeWorkflowExecution(ctx, workflowID, "")
	if err != nil {
		log.Printf("Failed to describe workflow %s: %v", workflowID, err)
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	status := map[string]interface{}{
		"workflowId":    workflowID,
		"status":        resp.WorkflowExecutionInfo.Status.String(),
		"startTime":     resp.WorkflowExecutionInfo.StartTime,
		"executionTime": resp.WorkflowExecutionInfo.ExecutionTime,
	}

	if resp.WorkflowExecutionInfo.CloseTime != nil {
		status["closeTime"] = resp.WorkflowExecutionInfo.CloseTime
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
