package request

type CreateSessionRequest struct {
	PaymentMethodId    float64 `json:"payment_method_id"`
	CustomerIdentifier string  `json:"customer_id"`
}
