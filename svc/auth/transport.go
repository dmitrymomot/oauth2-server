package auth

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/dmitrymomot/oauth2-server/internal/binder"
	"github.com/dmitrymomot/oauth2-server/internal/session"
	"github.com/dmitrymomot/oauth2-server/internal/utils"
	"github.com/dmitrymomot/oauth2-server/internal/validator"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func MakeHTTPHandler(srv Service, oauth2AuthURI string) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/login", httpLoginHandler(srv, oauth2AuthURI))
	r.HandleFunc("/register", httpRegisterHandler(srv, oauth2AuthURI))
	r.HandleFunc("/password/recovery", httpPasswordRecoveryHandler(srv))
	r.HandleFunc("/password/reset", httpPasswordResetHandler(srv))

	return r
}

// loginRequest collects the request parameters for the Login method.
type loginRequest struct {
	Email    string `json:"email" validate:"required|email" label:"Email address"`
	Password string `json:"password" validate:"required" label:"Password"`
}

// httpLoginHandler handles login requests.
func httpLoginHandler(srv Service, oauth2AuthURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if session.IsLoggedIn(r, w) {
			returnURI := session.GetReturnURI(r, w, oauth2AuthURI)
			w.Header().Set("Location", returnURI)
			w.WriteHeader(302)
			return
		}

		data := map[string]interface{}{
			"page_title": "Sign in",
		}

		if r.Method == http.MethodPost {
			payload := &loginRequest{}
			if err := binder.Bind(r, payload); err != nil {
				data["errors"] = []string{err.Error()}
				goview.Render(w, http.StatusOK, "login", data)
				return
			}
			data["form"] = payload

			if v := validator.ValidateStruct(payload); len(v) > 0 {
				data["validation"] = v
				goview.Render(w, http.StatusOK, "login", data)
				return
			}

			uid, err := srv.Login(r.Context(), payload.Email, payload.Password)
			if err != nil {
				data["errors"] = []string{err.Error()}
				goview.Render(w, http.StatusOK, "login", data)
				return
			}

			if err := session.StoreLoggedInUserID(r, w, uid.String()); err != nil {
				data["errors"] = []string{err.Error()}
				goview.Render(w, http.StatusOK, "login", data)
				return
			}

			returnURI := session.GetReturnURI(r, w, oauth2AuthURI)
			w.Header().Set("Location", returnURI)
			w.WriteHeader(302)
			return
		}

		goview.Render(w, http.StatusOK, "login", data)
	}
}

// registerRequest collects the request parameters for the Register method.
type registerRequest struct {
	Email                string `json:"email" form:"email" validate:"required|email|realEmail" filter:"trim|lower|escapeJs|escapeHtml|sanitizeEmail" label:"Email"`
	Password             string `json:"password" form:"password" validate:"required|minLen:8|maxLen:100" label:"Password"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" validate:"requiredWith:Password|eqField:Password" label:"Password confirmation" message:"Password confirmation must match password"`
	Terms                bool   `json:"terms" form:"terms" validate:"required|bool" label:"Terms of service"`
}

// httpRegisterHandler handles register requests.
func httpRegisterHandler(srv Service, oauth2AuthURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if session.IsLoggedIn(r, w) {
			returnURI := session.GetReturnURI(r, w, oauth2AuthURI)
			w.Header().Set("Location", returnURI)
			w.WriteHeader(302)
			return
		}

		data := map[string]interface{}{
			"page_title": "Registration",
		}

		if r.Method == http.MethodPost {
			payload := &registerRequest{}
			if err := binder.Bind(r, payload); err != nil {
				data["errors"] = []string{err.Error()}
				goview.Render(w, http.StatusOK, "register", data)
				return
			}
			data["form"] = payload

			if v := validator.ValidateStruct(payload); len(v) > 0 {
				data["validation"] = v
				goview.Render(w, http.StatusOK, "register", data)
				return
			}

			if _, err := srv.Register(r.Context(), payload.Email, payload.Password); err != nil {
				if errors.Is(err, ErrEmailTaken) {
					data["validation"] = url.Values{
						"email": []string{err.Error()},
					}
				} else {
					data["errors"] = []string{err.Error()}
				}
				goview.Render(w, http.StatusOK, "register", data)
				return
			}

			w.Header().Set("Location", utils.AddQueryParams("/verification", map[string]interface{}{
				"email": payload.Email,
			}))
			w.WriteHeader(302)
			return
		}

		goview.Render(w, http.StatusOK, "register", data)
	}
}

// passwordRecoveryRequest collects the request parameters for the PasswordRecovery method.
type passwordRecoveryRequest struct {
	Email string `json:"email" validate:"required|email" label:"Email address"`
}

// httpPasswordRecoveryHandler handles password recovery requests.
func httpPasswordRecoveryHandler(srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]interface{}{
			"page_title": "Password recovery",
		}

		if r.Method == http.MethodPost {
			payload := &passwordRecoveryRequest{}
			if err := binder.Bind(r, payload); err != nil {
				data["errors"] = []string{err.Error()}
				goview.Render(w, http.StatusOK, "password_recovery", data)
				return
			}
			data["form"] = payload

			if v := validator.ValidateStruct(payload); len(v) > 0 {
				data["validation"] = v
				goview.Render(w, http.StatusOK, "password_recovery", data)
				return
			}

			if err := srv.PasswordRecovery(r.Context(), payload.Email); err != nil {
				data["validation"] = url.Values{
					"email": []string{err.Error()},
				}
				goview.Render(w, http.StatusOK, "password_recovery", data)
				return
			}

			data["page_title"] = "Password recovery email sent"
			goview.Render(w, http.StatusOK, "password_recovery_sent", data)
			return
		}

		goview.Render(w, http.StatusOK, "password_recovery", data)
	}
}

// passwordResetRequest collects the request parameters for the PasswordReset method.
type passwordResetRequest struct {
	Email                string `json:"email" validate:"-" label:"Email address"`
	OTP                  string `json:"otp" validate:"-" label:"OTP code"`
	Password             string `json:"password" validate:"required|minLen:8|maxLen:100" label:"Password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"requiredWith:Password|eqField:Password" label:"Password confirmation"`
}

// httpPasswordResetHandler handles password reset requests.
func httpPasswordResetHandler(srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]interface{}{
			"page_title": "Password reset",
		}

		payload := &passwordResetRequest{}
		if err := binder.Bind(r, payload); err != nil {
			data["errors"] = []string{err.Error()}
			goview.Render(w, http.StatusOK, "password_recovery", data)
			return
		}
		data["form"] = payload

		if r.Method == http.MethodPost {
			if v := validator.ValidateStruct(payload); len(v) > 0 {
				data["validation"] = v
				goview.Render(w, http.StatusOK, "password_recovery", data)
				return
			}

			if err := srv.PasswordReset(
				r.Context(),
				payload.Email,
				payload.OTP,
				payload.Password,
			); err != nil {
				data["errors"] = []string{err.Error()}
				goview.Render(w, http.StatusOK, "password_recovery", data)

			}

			data["page_title"] = "Password has been changed"
			goview.Render(w, http.StatusOK, "password_reset_success", data)
			return
		}

		data["page_title"] = "Password reset"
		goview.Render(w, http.StatusOK, "password_recovery", data)
	}
}
