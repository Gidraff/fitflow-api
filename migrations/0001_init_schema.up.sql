-- 1. Setup Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Core Tables
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trainer_id TEXT NOT NULL, -- Maps to claims.sub from Keycloak
    name TEXT NOT NULL,
    goal TEXT,
    profile TEXT, -- General metadata
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'Unnamed Program',
    duration_weeks INT DEFAULT 8,
    is_active BOOLEAN DEFAULT true,
    plan JSONB NOT NULL, -- Flexible exercise structure (sets, reps, etc.)
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    weight NUMERIC(5,2), -- Supports decimal precision (e.g., 92.5kg)
    mood TEXT,
    energy_level INT CHECK (energy_level BETWEEN 1 AND 5),
    athlete_rating INT CHECK (athlete_rating BETWEEN 1 AND 5),
    performance_rating INT CHECK (performance_rating BETWEEN 1 AND 5),
    athlete_notes TEXT,
    trainer_feedback TEXT,
    commentary TEXT, -- Original field from v1
    created_at TIMESTAMPTZ DEFAULT now()
);

-- 3. Indexes for Performance
CREATE INDEX idx_clients_trainer ON clients(trainer_id);
CREATE INDEX idx_sessions_client ON sessions(client_id);
CREATE INDEX idx_workouts_client ON workouts(client_id);