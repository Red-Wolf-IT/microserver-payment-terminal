package storage_test

import (
	"testing"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/model"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/storage"
)

func TestPaymentStorage_CreateAndGetByID(t *testing.T) {
	store := storage.NewPaymentStorage()

	paymentID := "test-payment-id"

	createdPayment, err := store.Create(paymentID)
	if err != nil {
		t.Fatalf("failed to create payment: %v", err)
	}

	if createdPayment.PaymentID != paymentID {
		t.Errorf(
			"expected payment id %s, got %s",
			paymentID,
			createdPayment.PaymentID,
		)
	}

	if createdPayment.Status != model.PaymentStatusWaiting {
		t.Errorf(
			"expected status %s, got %s",
			model.PaymentStatusWaiting,
			createdPayment.Status,
		)
	}

	gotPayment, err := store.GetByID(paymentID)
	if err != nil {
		t.Fatalf("failed to get payment: %v", err)
	}

	if gotPayment.PaymentID != createdPayment.PaymentID {
		t.Errorf(
			"expected payment id %s, got %s",
			createdPayment.PaymentID,
			gotPayment.PaymentID,
		)
	}

	if gotPayment.Status != createdPayment.Status {
		t.Errorf(
			"expected status %s, got %s",
			createdPayment.Status,
			gotPayment.Status,
		)
	}
}

func TestPaymentStorage_Exists(t *testing.T) {
	store := storage.NewPaymentStorage()

	paymentID := "test-payment-id"

	_, err := store.Create(paymentID)
	if err != nil {
		t.Fatalf("failed to create payment: %v", err)
	}

	if !store.Exists(paymentID) {
		t.Errorf("expected payment to exist")
	}

	if store.Exists("wrong-id") {
		t.Errorf("payment with wrong id should not exist")
	}
}