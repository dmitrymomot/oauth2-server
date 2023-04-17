package session

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-session/session/v3"
)

const (
	// ReturnURIKey is the key used to store the return URI in the session.
	ReturnURIKey = "return_uri"
	// RedirectDataKey is the key used to store the redirect URI in the session.
	RedirectDataKey = "redirect_data"
	// LoggedInUserIDKey is the key used to store the logged in user ID in the session.
	LoggedInUserIDKey = "logged_in_user_id"
)

// StoreReturnURI stores the return URI in the session.
func StoreReturnURI(r *http.Request, w http.ResponseWriter, returnURI string) error {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return fmt.Errorf("session start: %w", err)
	}

	store.Set(ReturnURIKey, returnURI)
	if err := store.Save(); err != nil {
		return fmt.Errorf("session save: %w", err)
	}

	return nil
}

// StoreCurrentURI stores the current URI in the session as the return URI.
func StoreCurrentURI(r *http.Request, w http.ResponseWriter) error {
	return StoreReturnURI(r, w, r.URL.String())
}

// GetReturnURI gets the return URI from the session.
func GetReturnURI(r *http.Request, w http.ResponseWriter, fallbackURI string) string {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return fallbackURI
	}

	uri, ok := store.Get(ReturnURIKey)
	if !ok {
		return fallbackURI
	}

	result, ok := uri.(string)
	if !ok {
		return fallbackURI
	}

	// Delete the return URI from the session after it has been retrieved.
	store.Delete(ReturnURIKey)
	store.Save()

	return result
}

// StoreRedirectData stores the redirect URI in the session with request payload.
func StoreRedirectData(r *http.Request, w http.ResponseWriter) error {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return fmt.Errorf("session start: %w", err)
	}

	var form url.Values
	if r.Method == http.MethodPost {
		if r.Form == nil {
			r.ParseForm()
		}
		form = r.Form
	} else {
		form = r.URL.Query()
	}

	store.Set(RedirectDataKey, form)
	if err := store.Save(); err != nil {
		return fmt.Errorf("session save: %w", err)
	}

	return nil
}

// GetRedirectData gets the redirect URI from the session with request payload.
func GetRedirectData(r *http.Request, w http.ResponseWriter) url.Values {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return nil
	}

	payload, ok := store.Get(RedirectDataKey)
	if !ok {
		return nil
	}

	result, ok := payload.(url.Values)
	if !ok {
		return nil
	}

	// Delete the redirect URI from the session after it has been retrieved.
	store.Delete(RedirectDataKey)
	store.Save()

	return result
}

// StoreLoggedInUserID stores the logged in user ID in the session.
func StoreLoggedInUserID(r *http.Request, w http.ResponseWriter, userID string) error {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return fmt.Errorf("session start: %w", err)
	}

	store.Set(LoggedInUserIDKey, userID)
	if err := store.Save(); err != nil {
		return fmt.Errorf("session save: %w", err)
	}

	return nil
}

// GetLoggedInUserID gets the logged in user ID from the session.
func GetLoggedInUserID(r *http.Request, w http.ResponseWriter) (string, bool) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return "", false
	}

	userID, ok := store.Get(LoggedInUserIDKey)
	if !ok || userID == nil {
		return "", false
	}

	result, ok := userID.(string)
	if !ok || result == "" {
		return "", false
	}

	return result, true
}

// IsLoggedIn checks if the user is logged in.
func IsLoggedIn(r *http.Request, w http.ResponseWriter) bool {
	_, ok := GetLoggedInUserID(r, w)
	return ok
}

// Logout logs the user out.
func Logout(r *http.Request, w http.ResponseWriter) error {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return fmt.Errorf("session start: %w", err)
	}

	if err := session.Destroy(r.Context(), w, r); err != nil {
		return fmt.Errorf("session destroy: %w", err)
	}

	if err := store.Flush(); err != nil {
		return fmt.Errorf("session flush: %w", err)
	}

	log.Println("logged out")

	return nil
}
