package storage

import (
	"errors"
	"sync"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/model"
)

var (
	ErrPaymentInProgress = errors.New("payment already in progress")
	ErrPaymentNotFound   = errors.New("payment not found")
)

type PaymentStorage struct {
	mu      sync.RWMutex
	payment *model.Payment
}

func NewPaymentStorage() *PaymentStorage {
	return &PaymentStorage{}
}

func (s *PaymentStorage) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.payment != nil && s.payment.PaymentID == id
}

func (s *PaymentStorage) HasWaiting() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.payment != nil && s.payment.Status == model.PaymentStatusWaiting
}

func (s *PaymentStorage) Create(id string) (*model.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.payment != nil && s.payment.Status == model.PaymentStatusWaiting {
		return nil, ErrPaymentInProgress
	}

	payment := &model.Payment{
		PaymentID: id,
		Status:    model.PaymentStatusWaiting,
	}
	s.payment = payment

	paymentCopy := *payment
	return &paymentCopy, nil
}

func (s *PaymentStorage) GetByID(id string) (*model.Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.payment == nil || s.payment.PaymentID != id {
		return nil, ErrPaymentNotFound
	}

	paymentCopy := *s.payment
	return &paymentCopy, nil
}

func (s *PaymentStorage) UpdateStatus(id string, status model.PaymentStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.payment == nil || s.payment.PaymentID != id {
		return ErrPaymentNotFound
	}

	s.payment.Status = status
	return nil
}