package model

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type PaymentStatus string

const (
	PaymentStatusWaiting = "WAITING"
	PaymentStatusSuccess = "SUCCESS"
	PaymentStatusTimeout = "TIMEOUT"
)

type Payment struct {
	PaymentID  string         `json:"payment_id"`
	Status     PaymentStatus  `json:"status"`
}

func GeneratePaymentID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Error("generate payment id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func NewUniquePayment() (exists func(id string) bool) (string, error) {
	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		if, err := GeneratePaymentID()
		if err != nil {
			return "", err
		}
		if !exists(id) {
			return id, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique payment ID")
}