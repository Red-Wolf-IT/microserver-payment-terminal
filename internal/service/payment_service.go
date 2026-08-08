package service

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
	"io"

	"github.com/oklog/ulid/v2"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/model"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/storage"
)

type PaymentService struct {
	storage *storage.PaymentStorage

	entropyMu sync.Mutex
	entropy   io.Reader
}

func NewPaymentService(st *storage.PaymentStorage) *PaymentService {
	return &PaymentService{
		storage: st,
		entropy: ulid.Monotonic(rand.Reader, 0),
	}
}

func (s *PaymentService) CreatePayment() (*model.Payment, error) {
	id, err := s.newUniquePaymentID()
	if err != nil {
		return nil, err
	}

	return s.storage.Create(id)
}

func (s *PaymentService) GetPayment(id string) (*model.Payment, error) {
	return s.storage.GetByID(id)
}

func (s *PaymentService) generatePaymentID() (string, error) {
	s.entropyMu.Lock()
	defer s.entropyMu.Unlock()

	id, err := ulid.New(ulid.Timestamp(time.Now()), s.entropy)
	if err != nil {
		return "", fmt.Errorf("generate payment id: %w", err)
	}
	return id.String(), nil
}

func (s *PaymentService) newUniquePaymentID() (string, error) {
	const maxAttempts = 5

	for range maxAttempts {
		id, err := s.generatePaymentID()
		if err != nil {
			return "", err
		}
		if !s.storage.Exists(id) {
			return id, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique payment id")
}
