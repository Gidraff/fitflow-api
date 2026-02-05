package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"fitflow-api/internal/auth"
	"fitflow-api/internal/handlers"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	allowedOriginsRaw := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsRaw == "" {
		allowedOriginsRaw = "http://localhost:5173"
	}
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")

	db, err := sqlx.Connect("postgres", dbURL)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	v := validator.New()

	issuer := os.Getenv("KEYCLOAK_ISSUER_URL")
	jwksURL := fmt.Sprintf("%s/protocol/openid-connect/certs", issuer)

	// 3. Router Setup
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Preflight caching (5 minutes)
	}))

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
