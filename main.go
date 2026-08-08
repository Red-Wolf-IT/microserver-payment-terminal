package main

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/service"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/storage"
)

func main() {
	paymentStorage := storage.NewPaymentStorage()
	paymentService := service.NewPaymentService(paymentStorage)

	payment, err := paymentService.CreatePayment()
	if err != nil {
		log.Fatal(err)
	}

	jsonPayment, err := json.MarshalIndent(payment, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonPayment))
}