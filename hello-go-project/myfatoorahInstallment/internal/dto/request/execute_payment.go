package request

type ExecuteRequest struct {
	SessionId         string  `json:"SessionId"` // This is the Session B from frontend
	OriginalSessionID string  `json:"OriginalSessionID"`
	InvoiceValue      float64 `json:"InvoiceValue"`
	RecurringModel    struct {
		RecurringType string `json:"RecurringType"`
		IntervalDays  int    `json:"IntervalDays"`
		Iteration     int    `json:"Iteration"`
		RetryCount    int    `json:"RetryCount"`
	} `json:"RecurringModel"`
}
