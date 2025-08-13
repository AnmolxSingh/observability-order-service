package handler

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "math/rand"
    "strconv"
)

// In-memory storage
var orders = make(map[string]map[string]interface{})

// CreateOrder handles new orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    id := strconv.Itoa(rand.Intn(100000))
    order := map[string]interface{}{
        "id": id,
        "status": "created",
    }
    orders[id] = order
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}

// GetOrder returns order by ID
func GetOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    order, exists := orders[id]
    if !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}
