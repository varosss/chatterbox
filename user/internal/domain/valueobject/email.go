package valueobject

import (
	"errors"
	"strings"
)

type Email string

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(strings.ToLower(value))

	if value == "" {
		return "", errors.New("email is empty")
	}
	if !strings.Contains(value, "@") {
		return "", errors.New("invalid email")
	}

	return Email(value), nil
}

func (email Email) String() string {
	return string(email)
}
