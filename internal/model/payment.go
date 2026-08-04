package model

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

type PaymentStatus string

const (
	PaymentStatusWaiting PaymentStatus = "WAITING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusTimeout PaymentStatus = "TIMEOUT"
)

type Payment struct {
	PaymentID string        `json:"payment_id"`
	Status    PaymentStatus `json:"status"`
}

var (
	entropyMu sync.Mutex
	entropy   = ulid.Monotonic(rand.Reader, 0)
)

func GeneratePaymentID() (string, error) {
	entropyMu.Lock()
	defer entropyMu.Unlock()

	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return "", fmt.Errorf("generate payment id: %w", err)
	}
	return id.String(), nil
}

func NewUniquePayment(exists func(id string) bool) (string, error) {
	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		id, err := GeneratePaymentID()
		if err != nil {
			return "", err
		}
		if !exists(id) {
			return id, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique payment ID")
}