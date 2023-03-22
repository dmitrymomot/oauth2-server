package mailer

// Predefined task types.
const (
	SendConfirmationEmailTask     = "send_confirmation_email"
	SendPasswordRecoveryEmailTask = "send_password_recovery_email"
	SendDestroyProfileEmailTask   = "send_destroy_profile_email"
)

type (
	// Payload for sending confirmation email to user:
	// - email confirmation
	// - password reset
	// - destroy profile
	ConfirmationEmailPayload struct {
		UserID string `json:"user_id,omitempty"`
		Email  string `json:"email,omitempty"`
		OTP    string `json:"otp,omitempty"`
	}
)
