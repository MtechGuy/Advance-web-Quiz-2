package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (a *applicationDependencies) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the start time
		startTime := time.Now()

		log.Printf("Request received: Method=%s, URL=%s,", r.Method, r.RequestURI)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the completion of the request
		duration := time.Since(startTime)
		log.Printf("Request finished: Duration=%s", duration)
	})
}

// The AuthMiddleware function
func (a *applicationDependencies) AuthMiddleware(next http.Handler) http.Handler {
	username := []byte(a.Username)
	password := []byte(a.Password)
	realm := `Basic realm="Restricted"`
	if a.Realm != "" {
		realm = fmt.Sprintf(`Basic realm="%v"`, a.Realm)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), username) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), password) != 1 {

			w.Header().Set("WWW-Authenticate", realm)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized.\n"))
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (a *applicationDependencies) handleErrorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		// call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func (a *applicationDependencies) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}
