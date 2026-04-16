package request

import "github.com/google/uuid"

type CreateSessionRequest struct {
	PaymentMethodId    float64   `json:"payment_method_id"`
	CustomerIdentifier uuid.UUID `json:"customer_id"`
}
