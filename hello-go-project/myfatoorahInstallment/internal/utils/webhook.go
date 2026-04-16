package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func ValidateMyFatoorahSignature(body []byte, secret string, receivedSignature string) bool {
	// 1. Create a new HMAC using SHA256 and your secret key
	h := hmac.New(sha256.New, []byte(secret))

	// 2. Write the raw body bytes into the hash
	h.Write(body)

	// 3. Get the calculated signature and encode to Base64
	calculatedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 4. Compare (Use constant time compare if you want to be extra secure)
	return calculatedSignature == receivedSignature
}
