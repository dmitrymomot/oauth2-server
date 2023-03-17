package oauth

import (
	"net/http"
)

// getScopeFromRequest get scope from request
func getScopeFromRequest(r *http.Request) string {
	scopes := r.URL.Query().Get("scope")
	if scopes == "" {
		scopes = r.FormValue("scope")
	}
	return scopes
}
