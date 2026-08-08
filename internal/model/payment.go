package model

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
