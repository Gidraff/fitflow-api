package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"fitflow-api/internal/auth"
	"fitflow-api/internal/handlers"

	"github.com/go-chi/chi/v5/middleware" // This provides the standard middleware
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

func main() {
	_ = godotenv.Load()

	// 1. Database Initialization
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	v := validator.New()

	// 2. Auth Configuration (Reuse your Keycloak settings)
	issuer := os.Getenv("KEYCLOAK_ISSUER_URL")
	jwksURL := fmt.Sprintf("%s/protocol/openid-connect/certs", issuer)

	// 3. Router Setup
	r := chi.NewRouter()

	// Base Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Inject App State into Handlers
	state := &handlers.AppState{DB: db, Validator: v}

	// 4. Route Definitions
	r.Route("/api/v1", func(r chi.Router) {
		// Public Routes (if any)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware(jwksURL, issuer))

			// Client Endpoints
			r.Route("/clients", func(r chi.Router) {
				r.Get("/", state.ListClients)
				r.Post("/", state.CreateClient)
				r.Route("/{client_id}", func(r chi.Router) {
					r.Get("/", state.GetClientByID)
					r.Put("/", state.UpdateClient)
					r.Delete("/", state.DeleteClient)
					r.Get("/sessions", state.GetClientSessions)
				})
			})

			// Session Endpoints
			r.Post("/sessions", state.LogWorkoutSession)
			r.Put("/sessions/{session_id}/feedback", state.AddSessionFeedback)
		})
	})

	fmt.Println("Server starting on :8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
