package models

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrNotFound = errors.New("email not found")

type Email struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

func (e *Email) Validate() error {
	e.Address = strings.TrimSpace(e.Address)
	if e.Address == "" {
		return errors.New("address is required")
	}
	if _, err := mail.ParseAddress(e.Address); err != nil {
		return errors.New("address is not a valid email format")
	}
	return nil
}
