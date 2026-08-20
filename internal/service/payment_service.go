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
	storage  *storage.PaymentStorage
	timeout  time.Duration

	entropyMu sync.Mutex
	entropy   io.Reader

	timerMu sync.Mutex
	timer   *time.Timer
}

func NewPaymentService(st *storage.PaymentStorage, timeout time.Duration) *PaymentService {
	return &PaymentService{
		storage: st,
		timeout: timeout,
		entropy: ulid.Monotonic(rand.Reader, 0),
	}
}

func (s *PaymentService) CreatePayment() (*model.Payment, error) {
	id, err := s.newUniquePaymentID()
	if err != nil {
		return nil, err
	}
	payment, err := s.storage.Create(id)
	if err != nil {
		return nil, err
	}
	s.startTimeout(id)
	return payment, nil
}

func (s *PaymentService) startTimeout(id string) {
	s.timerMu.Lock()
	defer s.timerMu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
	}
	s.timer = time.AfterFunc(s.timeout, func() {
		_ = s.storage.Timeout(id)
	})
}
func (s *PaymentService) stopTimeout() {
	s.timerMu.Lock()
	defer s.timerMu.Unlock()
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
}
func (s *PaymentService) ConfirmPayment(id string) error {
	if err := s.storage.UpdateStatus(id, model.PaymentStatusSuccess); err != nil {
		return err
	}
	s.stopTimeout()
	return nil
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
