package utils

import "strings"

// MaskNationalID masks a national ID, showing only the last 4 characters.
func MaskNationalID(id string) string {
	if len(id) <= 4 {
		return strings.Repeat("*", len(id))
	}
	return strings.Repeat("*", len(id)-4) + id[len(id)-4:]
}

// MaskEmail masks an email address, showing only the first 2 characters and domain.
func MaskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email
	}
	local := parts[0]
	if len(local) <= 2 {
		return local + "***@" + parts[1]
	}
	return local[:2] + "***@" + parts[1]
}

// MaskPhone masks a phone number, showing only the last 4 digits.
func MaskPhone(phone string) string {
	if len(phone) <= 4 {
		return strings.Repeat("*", len(phone))
	}
	return strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
}
