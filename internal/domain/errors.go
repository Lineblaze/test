package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}
