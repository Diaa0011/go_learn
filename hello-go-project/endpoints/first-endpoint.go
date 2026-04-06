package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Define the request structure based on your JSON
type MyFatoorahRequest struct {
	PaymentMethod      string `json:"PaymentMethod"`
	Order              Order  `json:"Order"`
	PaymentExpiry      string `json:"PaymentExpiry"`
	NotificationOption string `json:"NotificationOption"`
	Language           string `json:"Language"`
}

type Order struct {
	Amount   string `json:"Amount"`
	Currency string `json:"Currency"`
}

func main() {
	apiUrl := "https://apitest.myfatoorah.com/v3/payments"
	token := "SK_KWT_vVZlnnAqu8jRByOWaRPNId4ShzEDNt256dvnjebuyzo52dXjAfRx2ixW5umjWSUx"

	// Construct the payload
	payload := MyFatoorahRequest{
		PaymentMethod: "CARD",
		Order: Order{
			Amount:   "1",
			Currency: "KWD",
		},
		PaymentExpiry:      "2026-04-13T08:30:00Z",
		NotificationOption: "LNK",
		Language:           "EN",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Create request
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print response
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status Code:", resp.StatusCode)

	// Pretty print the JSON response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "    "); err == nil {
		fmt.Println("Response Body:\n", prettyJSON.String())
	} else {
		fmt.Println("Raw Response:", string(body))
	}
}
