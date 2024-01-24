package middleware

import (
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
		_, err := auth.decoder.Decode(token)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
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
