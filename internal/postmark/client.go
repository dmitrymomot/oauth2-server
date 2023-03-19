package postmark_client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/keighl/postmark"
)

// Predefined email templates
var (
	VerificationCodeTmpl = "verification_code"
	PasswordResetTmpl    = "password_reset"
	DestroyUserCodeTmpl  = "destroy_account"
)

type (
	// Postmark service wrapper
	Client struct {
		client postmarkClient
		config Config
	}

	// Config struct
	Config struct {
		ProductName  string
		ProductURL   string
		SupportEmail string
		CompanyName  string
		FromEmail    string
		FromName     string

		// Base URLs for email actions
		VerificationCodeURL string
		PasswordResetURL    string
		DestroyUserCodeURL  string
	}

	postmarkClient interface {
		SendTemplatedEmail(email postmark.TemplatedEmail) (postmark.EmailResponse, error)
	}
)

// New func is a factory function,
// returns a new instance of the Client interface implementation
func New(postmark postmarkClient, conf Config) *Client {
	return &Client{
		client: postmark,
		config: conf,
	}
}

func (c *Client) SendVerificationCode(ctx context.Context, uid, email, otp string) error {
	actionURL, err := url.Parse(c.config.VerificationCodeURL)
	if err != nil {
		return fmt.Errorf("could not parse action url: %w", err)
	}
	actionURL.RawQuery = url.Values{
		"uid":   {uid},
		"otp":   {otp},
		"email": {email},
	}.Encode()

	return c.send(
		VerificationCodeTmpl,
		"verification_code",
		email,
		map[string]interface{}{
			"otp":        otp,
			"action_url": actionURL.String(),
		},
	)
}

func (c *Client) SendResetPasswordCode(ctx context.Context, uid, email, otp string) error {
	actionURL, err := url.Parse(c.config.PasswordResetURL)
	if err != nil {
		return fmt.Errorf("could not parse action url: %w", err)
	}
	actionURL.RawQuery = url.Values{
		"uid":   {uid},
		"otp":   {otp},
		"email": {email},
	}.Encode()

	return c.send(
		PasswordResetTmpl,
		"reset_password",
		email,
		map[string]interface{}{
			"otp":        otp,
			"action_url": actionURL.String(),
		},
	)
}

func (c *Client) SendDestroyProfileCode(ctx context.Context, uid, email, otp string) error {
	actionURL, err := url.Parse(c.config.DestroyUserCodeURL)
	if err != nil {
		return fmt.Errorf("could not parse action url: %w", err)
	}
	actionURL.RawQuery = url.Values{
		"uid":   {uid},
		"otp":   {otp},
		"email": {email},
	}.Encode()

	return c.send(
		DestroyUserCodeTmpl,
		"destroy_account",
		email,
		map[string]interface{}{
			"otp":        otp,
			"action_url": actionURL.String(),
		},
	)
}

// send email
func (c *Client) send(tpl, tag, email string, data map[string]interface{}) error {
	// Default model data
	payload := map[string]interface{}{
		"product_url":   c.config.ProductURL,
		"product_name":  c.config.ProductName,
		"company_name":  c.config.CompanyName,
		"email":         email,
		"support_email": c.config.SupportEmail,
	}

	// Merge custom data with default fields
	for k, v := range data {
		payload[k] = v
	}

	if _, err := c.client.SendTemplatedEmail(postmark.TemplatedEmail{
		TemplateAlias: tpl,
		InlineCss:     true,
		TrackOpens:    true,
		From:          c.config.FromEmail,
		To:            email,
		Tag:           tag,
		ReplyTo:       c.config.SupportEmail,
		TemplateModel: payload,
	}); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
