package response

type CreateSessionResponse struct {
	IsSuccess bool `json:"IsSuccess"`
	Data      struct {
		SessionId string `json:"SessionId"`
		Country   string `json:"CountryCode"`
	} `json:"Data"`
	Message string `json:"Message"`
}
