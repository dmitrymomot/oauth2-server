package mdw

import (
	"net/http"

	"github.com/dmitrymomot/oauth2-server/internal/session"
)

// Check if user is not authenticated and redirect to home page.
// If user is not authenticated, continue to next handler. Otherwise, redirect to home page.
func NotAuthOnly(homeURI string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session.IsLoggedIn(r, w) {
				returnURI := session.GetReturnURI(r, w, homeURI)
				http.Redirect(w, r, returnURI, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
