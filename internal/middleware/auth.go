package middleware

import (
	"context"
	"mfv-challenge/internal/constants"
	"net/http"
	"strings"
)

type Authenticator struct {
	decoder   TokenDecoder
	enabled   bool
	whitelist map[string]bool
}

type TokenDecoder interface {
	Decode(tokenString string) (int64, error)
}

type contextKey string

const authenticatedUserKey contextKey = "X-UserID"

// Middleware function, which will be called for each request
func (auth *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.enabled {
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := auth.whitelist[r.URL.String()]; ok {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get(constants.AuthorizationHeader)
		token = strings.ReplaceAll(token, constants.AuthorizationKey, "")
		userID, err := auth.decoder.Decode(token)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		//create a new request context containing the authenticated user
		ctxWithUser := context.WithValue(r.Context(), authenticatedUserKey, userID)
		//create a new request using that new context
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)
	})
}

func NewAuth(decoder TokenDecoder, whitelist []string, enabled bool) *Authenticator {
	m := map[string]bool{}
	for _, url := range whitelist {
		m[url] = true
	}
	return &Authenticator{
		decoder:   decoder,
		whitelist: m,
		enabled:   enabled,
	}
}
