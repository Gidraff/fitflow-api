package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	// Update the import path below to the correct location of your models package.
	// For example, if your models are in /Users/gidraff/GoProjects/fitflow-api/internal/models, use:
	"fitflow-api/internal/models"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwksURL string, issuer string) func(http.Handler) http.Handler {
	var jwks *keyfunc.JWKS
	var err error

	// 1. Retry loop: Keycloak takes time to boot in Docker
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		options := keyfunc.Options{
			RefreshInterval: time.Hour,
		}

		jwks, err = keyfunc.Get(jwksURL, options)
		if err == nil {
			break // Success!
		}

		log.Printf("Keycloak not ready yet (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second) // Wait before trying again
	}

	// If we still fail after all retries, then we stop the app
	if err != nil {
		log.Fatalf("Failed to get JWKS from Keycloak after %d retries: %v", maxRetries, err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ... (Rest of your existing token extraction and validation logic)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			var claims models.Claims
			token, err := jwt.ParseWithClaims(tokenString, &claims, jwks.Keyfunc)

			if err != nil || !token.Valid {
				log.Printf("JWT Error: %v\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims.Issuer != issuer {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", &claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
