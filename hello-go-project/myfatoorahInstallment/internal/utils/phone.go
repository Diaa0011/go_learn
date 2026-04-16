package utils

import (
	"strings"
)

func SplitPhoneAndCode(fullPhone string) (countryCode string, localNumber string) {
	if len(fullPhone) == 0 {
		return "", ""
	}

	// Remove the '+' if present
	cleanPhone := fullPhone
	if fullPhone[0] == '+' {
		cleanPhone = fullPhone[1:]
	}

	// Try to match against our generic list
	for _, code := range CountryCodes {
		if strings.HasPrefix(cleanPhone, code) {
			return code, cleanPhone[len(code):]
		}
	}

	// Fallback: if no match found, return the whole thing as the number
	return "", cleanPhone
}
