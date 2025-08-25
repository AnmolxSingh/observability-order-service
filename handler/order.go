package handler

import (
	"bytes"
	"encoding/json"
	"io"
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
		Amount    int    `json:"amount"`
	}

	// Parse JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	request1, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/inventory", nil)
	// Inject tracing headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request1.Header))
	// Send request
	client := http.Client{}
	resp1, err := client.Do(request1)
	if err != nil {
		http.Error(w, "failed to contact inventory", 500)
		return
	}
	defer resp1.Body.Close()

	// 2. Call Payment Service
	paymentReq := map[string]interface{}{
		"orderId": "temp-order", // we’ll replace with actual ID later
		"user":    req.User,
		"amount":  req.Amount,
	}
	paymentReqBody, _ := json.Marshal(paymentReq)

	request2, _ := http.NewRequestWithContext(ctx, "POST", "http://localhost:8082/payment", bytes.NewBuffer(paymentReqBody))
	request2.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request2.Header))

	resp2, err := client.Do(request2)
	if err != nil {
		http.Error(w, "failed to contact payment", http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// Check HTTP status
	if resp2.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp2.Body)
		http.Error(w, "payment service error: "+string(bodyBytes), resp2.StatusCode)
		return
	}

	// Parse Payment Response
	var paymentResp map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&paymentResp); err != nil {
		http.Error(w, "invalid payment response", http.StatusInternalServerError)
		return
	}

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
