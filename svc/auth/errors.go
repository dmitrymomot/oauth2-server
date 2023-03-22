package auth

import "errors"

// Predefined errors
var (
	ErrEmailTaken                 = errors.New("Email is already taken")
	ErrInvalidCredentials         = errors.New("Invalid credentials. Please check your email and password and try again.")
	ErrUserNotFound               = errors.New("User not found")
	ErrInvalidVerificationRequest = errors.New("Invalid verification request")
	ErrInvalidVerificationCode    = errors.New("Invalid verification code")
	ErrVerificationCodeExpired    = errors.New("Verification code expired")
	ErrUserNotVerified            = errors.New("User not verified")
	ErrUserAlreadyVerified        = errors.New("User already verified")
)
