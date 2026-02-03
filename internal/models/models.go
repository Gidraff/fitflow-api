package models

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	Email   string    `db:"email" json:"email"` // Added
	Goal    *string   `db:"goal" json:"goal"`
	Profile *string   `db:"profile" json:"profile"`
}

type CreateClientRequest struct {
	Name    string  `json:"name" validate:"required,min=3"`
	Email   string  `json:"email" validate:"required,email"` // Added
	Goal    *string `json:"goal" validate:"required,max=100"`
	Profile *string `json:"profile"`
}

type Session struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	ClientID          uuid.UUID  `db:"client_id" json:"client_id"`
	WorkoutID         *uuid.UUID `db:"workout_id" json:"workout_id"`
	Date              time.Time  `db:"date" json:"date"`
	Weight            *float64   `db:"weight" json:"weight"`
	Mood              *string    `db:"mood" json:"mood"`
	EnergyLevel       *int       `db:"energy_level" json:"energy_level"`
	AthleteRating     *int       `db:"athlete_rating" json:"athlete_rating"`
	AthleteNotes      *string    `db:"athlete_notes" json:"athlete_notes"`
	TrainerFeedback   *string    `db:"trainer_feedback" json:"trainer_feedback"`
	PerformanceRating *int       `db:"performance_rating" json:"performance_rating"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
}
