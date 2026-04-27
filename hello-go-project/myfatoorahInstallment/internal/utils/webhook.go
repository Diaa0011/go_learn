package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

func BuildSignture(payload map[string]interface{}) string {
	data, ok := payload["Data"].(map[string]interface{})
	if !ok {
		return "no"
	}

	invoice := data["Invoice"].(map[string]interface{})
	transaction := data["Transaction"].(map[string]interface{})

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Invoice.Id=%v,", getString(invoice["Id"])))
	sb.WriteString(fmt.Sprintf("Invoice.Status=%v,", getString(invoice["Status"])))

	sb.WriteString(fmt.Sprintf("Transaction.Status=%v,", getString(transaction["Status"])))
	sb.WriteString(fmt.Sprintf("Transaction.PaymentId=%v,", getString(transaction["PaymentId"])))

	sb.WriteString(fmt.Sprintf("Invoice.ExternalIdentifier=%v", getString(invoice["ExternalIdentifier"])))

	return sb.String()
}

func ValidateMyFatoorahSignature(body []byte, secret string, receivedSignature string) bool {
	// 1. Create a new HMAC using SHA256 and your secret key
	h := hmac.New(sha256.New, []byte(secret))

	// 2. Write the raw body bytes into the hash
	h.Write(body)

	// 3. Get the calculated signature and encode to Base64
	calculatedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	log.Printf("Calculated Signature: %s", calculatedSignature)
	log.Printf("Received Signature: %s", receivedSignature)

	// 4. Compare (Use constant time compare if you want to be extra secure)
	return calculatedSignature == receivedSignature
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
