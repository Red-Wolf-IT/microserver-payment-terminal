package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/service"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/storage"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/handler"
)

func main() {
	paymentStorage := storage.NewPaymentStorage()
	paymentService := service.NewPaymentService(paymentStorage, 60*time.Second)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/payments", paymentHandler.CreatePayment)
	mux.HandleFunc("GET /api/v1/payments", paymentHandler.GetPayment)
	mux.HandleFunc("POST /api/v1/payments/confirm", paymentHandler.ConfirmPayment)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}