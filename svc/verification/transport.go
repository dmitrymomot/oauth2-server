package verification

import (
	"fmt"
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/binder"
	"github.com/dmitrymomot/oauth2-server/internal/validator"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func MakeHTTPHandler(srv Service) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/", httpVerificationHandler(srv))
	r.Post("/resend", httpResendHandler(srv))
	r.Get("/link", httpVerificationLinkHandler(srv))

	return r
}

// httpVerificationHandlerRequest is the request payload for the verification handler.
type httpVerificationHandlerRequest struct {
	Email string `json:"email" validate:"required|email" filter:"trim|lower|escapeJs|escapeHtml|sanitizeEmail" label:"Email"`
	OTP   string `json:"otp" validate:"required" filter:"trim|escapeJs|escapeHtml" label:"OTP"`
}

// httpVerificationHandler handles verification requests.
func httpVerificationHandler(srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]interface{}{
			"page_title": "Email Verification",
		}

		payload := httpVerificationHandlerRequest{}
		if err := binder.Bind(r, &payload); err != nil {
			data["errors"] = []string{err.Error()}
			goview.Render(w, http.StatusOK, "verification", data)
			return
		}
		data["form"] = payload

		goview.Render(w, http.StatusOK, "verification", data)
	}
}

// httpResendHandlerRequest is the request payload for the resend handler.
type httpResendHandlerRequest struct {
	Email string `json:"email" validate:"required|email" filter:"trim|lower|escapeJs|escapeHtml|sanitizeEmail" label:"Email"`
}

// httpResendHandler handles resend requests.
func httpResendHandler(srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]interface{}{
			"page_title": "Email Verification",
		}

		payload := &httpResendHandlerRequest{}
		if err := binder.Bind(r, payload); err != nil {
			data["errors"] = []string{err.Error()}
			goview.Render(w, http.StatusOK, "verification", data)
			return
		}
		data["form"] = payload

		if v := validator.ValidateStruct(payload); len(v) > 0 {
			data["validation"] = v
			goview.Render(w, http.StatusOK, "verification", data)
			return
		}

		if err := srv.ResendEmail(r.Context(), payload.Email); err != nil {
			data["errors"] = []string{err.Error()}
		} else {
			data["success"] = []string{fmt.Sprintf("New verification email has been sent to %s", payload.Email)}
		}

		goview.Render(w, http.StatusOK, "verification", data)
		return
	}
}

// httpVerificationLinkHandlerRequest is the request payload for the verification handler.
type httpVerificationLinkHandlerRequest struct {
	Email string `json:"email" validate:"required|email" filter:"trim|lower|escapeJs|escapeHtml|sanitizeEmail" label:"Email"`
	OTP   string `json:"otp" validate:"required" filter:"trim|lower|escapeJs|escapeHtml" label:"OTP"`
}

// httpVerificationLinkHandler handles verification link requests.
func httpVerificationLinkHandler(srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]interface{}{
			"page_title": "Email Verification",
		}

		payload := &httpVerificationLinkHandlerRequest{}
		if err := binder.Bind(r, payload); err != nil {
			data["errors"] = []string{err.Error()}
			goview.Render(w, http.StatusOK, "verification", data)
			return
		}
		data["form"] = payload

		if v := validator.ValidateStruct(payload); len(v) > 0 {
			data["validation"] = v
			goview.Render(w, http.StatusOK, "verification", data)
			return
		}

		if err := srv.VerifyEmail(r.Context(), payload.Email, payload.OTP); err != nil {
			data["errors"] = []string{err.Error()}
			goview.Render(w, http.StatusOK, "verification", data)
			return
		}

		goview.Render(w, http.StatusOK, "verification_success", data)
	}
}
