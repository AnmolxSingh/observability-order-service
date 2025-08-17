package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// In-memory storage
var orders = make(map[string]map[string]interface{})

// CreateOrder handles new orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Optional: tracing span if using OpenTelemetry
	tr := otel.Tracer("order-service")
	ctx, span := tr.Start(r.Context(), "CreateOrder")
	defer span.End()

	// Define request structure
	var req struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
		User      string `json:"user"`
	}

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	request, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/check", nil)

	// Inject tracing headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	// Send request
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		http.Error(w, "failed to contact inventory", 500)
		return
	}
	defer resp.Body.Close()

	// Create order with input fields
	id := strconv.Itoa(rand.Intn(100000))
	order := map[string]interface{}{
		"id":        id,
		"productId": req.ProductID,
		"quantity":  req.Quantity,
		"user":      req.User,
		"status":    "created",
	}

	// Store in memory
	orders[id] = order

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetOrder returns order by ID
// func GetOrder(w http.ResponseWriter, r *http.Request) {
// 	tr := otel.Tracer("order-service")
// 	_, span := tr.Start(r.Context(), "GetOrder")
// 	vars := mux.Vars(r)
// 	id := vars["id"]
// 	order, exists := orders[id]
// 	if !exists {
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(order)
// }
