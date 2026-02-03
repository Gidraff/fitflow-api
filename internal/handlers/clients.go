package handlers

import (
	"encoding/json"
	"fitflow-api/internal/models"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type AppState struct {
	DB        *sqlx.DB
	Validator *validator.Validate // Add this
}

func (s *AppState) CreateClient(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	var payload models.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Validate the input (Checks name length, email format, etc.)
	if err := s.Validator.Struct(payload); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	var client models.Client
	// 2. Insert with the email field
	query := `INSERT INTO clients (trainer_id, name, email, goal, profile) 
              VALUES ($1, $2, $3, $4, $5) 
              RETURNING id, name, email, goal, profile`

	err := s.DB.Get(&client, query, claims.Sub, payload.Name, payload.Email, payload.Goal, payload.Profile)

	if err != nil {
		// 3. Handle the Duplicate Client error
		// We look for the name of the constraint we added: unique_trainer_client_email
		if strings.Contains(err.Error(), "unique_trainer_client_email") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "A client with this email already exists."})
			return
		}

		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func (s *AppState) GetClientByID(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "client_id")
	claims := r.Context().Value("claims").(*models.Claims)

	var client models.Client
	err := s.DB.Get(&client, "SELECT id, name, goal, profile FROM clients WHERE id = $1 AND trainer_id = $2", clientID, claims.Sub)
	if err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(client)
}

func (s *AppState) LogWorkoutSession(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Helper to extract values safely (mimics payload["client_id"].as_str())
	getString := func(key string) string {
		if val, ok := body[key].(string); ok {
			return val
		}
		return ""
	}

	query := `
		INSERT INTO sessions (client_id, workout_id, date, weight, mood, energy_level, athlete_rating, athlete_notes) 
		VALUES ($1, $2, CURRENT_DATE, $3, $4, $5, $6, $7)`

	_, err := s.DB.Exec(query,
		getString("client_id"),
		getString("workout_id"),
		body["weight"],
		getString("mood"),
		body["energy_level"],
		body["athlete_rating"],
		getString("athlete_notes"),
	)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// List all clients for the authenticated trainer
func (s *AppState) ListClients(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)
	var clients []models.Client

	err := s.DB.Select(&clients, "SELECT id, name, goal, profile FROM clients WHERE trainer_id = $1", claims.Sub)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// Update an existing client
func (s *AppState) UpdateClient(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "client_id")
	claims := r.Context().Value("claims").(*models.Claims)

	var payload models.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	var client models.Client
	query := `UPDATE clients SET name = $1, goal = $2, profile = $3 
	          WHERE id = $4 AND trainer_id = $5 
	          RETURNING id, name, goal, profile`

	err := s.DB.Get(&client, query, payload.Name, payload.Goal, payload.Profile, clientID, claims.Sub)
	if err != nil {
		http.Error(w, "Client not found or update failed", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(client)
}

// Delete a client
func (s *AppState) DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "client_id")
	claims := r.Context().Value("claims").(*models.Claims)

	result, err := s.DB.Exec("DELETE FROM clients WHERE id = $1 AND trainer_id = $2", clientID, claims.Sub)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetClientSessions fetches all workout sessions for a specific client
func (s *AppState) GetClientSessions(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "client_id")
	claims := r.Context().Value("claims").(*models.Claims)

	var sessions []models.Session
	query := `SELECT id, date, weight, mood, energy_level, athlete_rating, athlete_notes 
	          FROM sessions WHERE client_id = $1 AND trainer_id = $2`

	err := s.DB.Select(&sessions, query, clientID, claims.Sub)
	if err != nil {
		http.Error(w, "Could not fetch sessions", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessions)
}

// AddSessionFeedback allows a trainer to add notes to a specific session
func (s *AppState) AddSessionFeedback(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "session_id")
	var body struct {
		Feedback string `json:"feedback"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err := s.DB.Exec("UPDATE sessions SET trainer_notes = $1 WHERE id = $2", body.Feedback, sessionID)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
