package utils

import (
	"regexp"
	"strings"
)

var (
	kenyanPhoneRegex  = regexp.MustCompile(`^(?:\+254|254|0)?([17]\d{8})$`)
	emailRegex        = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	nationalIDRegex   = regexp.MustCompile(`^\d{7,8}$`)
)

func ValidateKenyanPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	return kenyanPhoneRegex.MatchString(phone)
}

func NormalizeKenyanPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	matches := kenyanPhoneRegex.FindStringSubmatch(phone)
	if len(matches) < 2 {
		return phone
	}
	return "+254" + matches[1]
}

func ValidateEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	return emailRegex.MatchString(email)
}

func ValidateNationalID(id string) bool {
	id = strings.TrimSpace(id)
	return nationalIDRegex.MatchString(id)
}
