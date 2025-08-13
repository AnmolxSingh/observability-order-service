package main

import (
    "context"
    "log"
    "net/http"
    "order-service/handler"
    "order-service/otel"

    "github.com/gorilla/mux"
)

func main() {
    // Setup OpenTelemetry
    tp, err := otel.InitTracer()
    if err != nil {
        log.Fatalf("failed to initialize tracer: %v", err)
    }
    defer tp.Shutdown(context.Background())

    r := mux.NewRouter()
    r.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
    r.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")

    log.Println("Order Service running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
