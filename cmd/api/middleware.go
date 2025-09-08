package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// BasicAuthMiddleware applies HTTP Basic auth to protected endpoints.
func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization is missing"))
				return
			}

			// Expect 'Basic <base64>'
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// Decode base64 credentials
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			// Verify username and password
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
