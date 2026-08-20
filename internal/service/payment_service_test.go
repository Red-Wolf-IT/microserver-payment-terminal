package service_test

import (
	"testing"
	"time"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/model"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/service"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/storage"
)
func TestCreatePayment_TimesOut(t *testing.T) {
	store := storage.NewPaymentStorage()
	svc := service.NewPaymentService(store, 20*time.Millisecond)
	payment, err := svc.CreatePayment()

	if err != nil {
		t.Fatalf("create: %v", err)
	}

	time.Sleep(40 * time.Millisecond)

	got, err := svc.GetPayment(payment.PaymentID)

	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got.Status != model.PaymentStatusTimeout {
		t.Fatalf("expected TIMEOUT, got %s", got.Status)
	}

	if err := svc.ConfirmPayment(payment.PaymentID); err == nil {
		t.Fatal("confirm after timeout should fail")
	}

	if _, err := svc.CreatePayment(); err != nil {
		t.Fatalf("next create should succeed: %v", err)
	}
}

func TestConfirmPayment_StopsTimeout(t *testing.T) {
	store := storage.NewPaymentStorage()
	svc := service.NewPaymentService(store, 50*time.Millisecond)
	payment, err := svc.CreatePayment()

	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := svc.ConfirmPayment(payment.PaymentID); err != nil {
		t.Fatalf("confirm: %v", err)
	}

	time.Sleep(80 * time.Millisecond)
	
	got, err := svc.GetPayment(payment.PaymentID)

	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got.Status != model.PaymentStatusSuccess {
		t.Fatalf("expected SUCCESS, got %s", got.Status)
	}
}