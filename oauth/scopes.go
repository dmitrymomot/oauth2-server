package oauth

import (
	"strings"

	scopes "github.com/SonicRoshan/scope"
)

// MatchScopesStrict verifies if the all scopes is allowed.
// It returns true if the scope is allowed, false otherwise.
func MatchScopesStrict(requiredScopes string, allowedScopes string) bool {
	allowedList := strings.Split(allowedScopes, " ")
	if len(allowedList) == 0 {
		return false
	}
	scopeList := strings.Split(requiredScopes, " ")
	if len(scopeList) == 0 {
		return false
	}
	for _, s := range scopeList {
		if !scopes.ScopeInAllowed(s, allowedList) {
			return false
		}
	}
	return true
}

// MatchScopes verifies if the scope is allowed.
// It returns true if the scope is allowed, false otherwise.
func MatchScopes(requiredScopes string, allowedScopes string) bool {
	allowedList := strings.Split(allowedScopes, " ")
	if len(allowedList) == 0 {
		return false
	}
	scopeList := strings.Split(requiredScopes, " ")
	if len(scopeList) == 0 {
		return false
	}
	for _, s := range scopeList {
		if scopes.ScopeInAllowed(s, allowedList) {
			return true
		}
	}
	return false
}

// MatchScopes verifies if the scope is allowed.
// It returns true if the scope is allowed, false otherwise.
func MatchScope(requiredScope string, allowedScopes string) bool {
	if allowedScopes == "" && requiredScope != "" {
		return false
	} else if requiredScope == "" {
		return true
	}

	allowedList := strings.Split(allowedScopes, " ")
	if allowedList == nil || len(allowedList) < 1 {
		return false
	}
	return scopes.ScopeInAllowed(requiredScope, allowedList)
}
