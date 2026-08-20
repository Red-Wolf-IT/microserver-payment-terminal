package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/service"
	"github.com/Red-Wolf-IT/microserver-payment-terminal/internal/storage"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	payment, err := h.paymentService.CreatePayment()
	if err != nil {
		if errors.Is(err, storage.ErrPaymentInProgress) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create payment"})
		return
	}

	writeJSON(w, http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "payment ID is required"})
		return
	}

	payment, err := h.paymentService.GetPayment(id)
	if err != nil {
		if errors.Is(err, storage.ErrPaymentNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get payment"})
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "payment ID is required"})
		return
	}

	if err := h.paymentService.ConfirmPayment(id); err != nil {
		if errors.Is(err, storage.ErrPaymentNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, storage.ErrPaymentNotActive) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to confirm payment"})
		return
	}

	payment, err := h.paymentService.GetPayment(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get payment"})
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}