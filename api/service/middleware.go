package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adrianosela/padl/api/auth"
)

// we need a type for context key
type ctxKey string

var (
	// AccessTokenClaimsKey is the key in the request
	// context object for access token claims
	AccessTokenClaimsKey = ctxKey("access-claims")
)

// Auth wraps an HTTP handler function
// and populates the access token claims object in the req ctx
func (s *Service) Auth(h http.HandlerFunc, allowedAuds ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from request header
		authorization := r.Header.Get("Authorization")
		tkStr := strings.TrimPrefix(authorization, "Bearer ")
		if authorization == tkStr {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "no access token in header")
			return
		}
		// validate token
		verifiedClaims, err := s.authenticator.ValidateJWT(tkStr, allowedAuds...)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "invalid access token")
			return
		}

		// run handler with token in context
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AccessTokenClaimsKey, verifiedClaims)))
	})
}

// GetClaims returns the claims in a context object
func GetClaims(r *http.Request) *auth.CustomClaims {
	return r.Context().Value(AccessTokenClaimsKey).(*auth.CustomClaims)
}
