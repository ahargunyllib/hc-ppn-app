package phoneutil

import "strings"

// NormalizeToE164 converts a phone number to E.164 format by adding the + prefix if not present.
// E.164 is the international standard for phone numbers: +[country code][number]
//
// Examples:
//   - "628123456789" -> "+628123456789"
//   - "+628123456789" -> "+628123456789"
//   - "" -> ""
func NormalizeToE164(phoneNumber string) string {
	trimmed := strings.TrimSpace(phoneNumber)

	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "+") {
		return trimmed
	}

	return "+" + trimmed
}
