package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.temporal.io/sdk/client"
)

type HealthHandlers struct {
}

func NewHealthHandlers(temporalClient client.Client) *HealthHandlers {
	return &HealthHandlers{}
}

func (h *HealthHandlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}
