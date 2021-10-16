package api

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

// AdminAuthorization is a default authorizor middleware to enforce resource access
func AdminAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if claims["admin"] == nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authorized, pass it through
		next.ServeHTTP(w, r)
	})
}
