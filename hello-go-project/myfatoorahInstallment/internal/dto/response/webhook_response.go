package response

type WebhookPayload struct {
	Event struct {
		Name string `json:"Name"`
	} `json:"Event"`
	Data struct {
		Invoice struct {
			Id                 string                 `json:"Id"`
			Status             string                 `json:"Status"`
			ExternalIdentifier string                 `json:"ExternalIdentifier"` // This is your Invoice UUID!
			MetaData           map[string]interface{} `json:"MetaData"`
		} `json:"Invoice"`
		Transaction struct {
			Id        string `json:"Id"`
			Status    string `json:"Status"` // "SUCCESS"
			PaymentId string `json:"PaymentId"`
		} `json:"Transaction"`
	} `json:"Data"`
	Amount struct {
		ValueInDisplayCurrency string `json:"ValueInDisplayCurrency"`
	} `json:"Amount"`
}
