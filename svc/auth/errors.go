package auth

import "errors"

// Predefined errors
var (
	ErrEmailTaken         = errors.New("Email is already taken")
	ErrInvalidCredentials = errors.New("Invalid credentials. Please check your email and password and try again.")
)
