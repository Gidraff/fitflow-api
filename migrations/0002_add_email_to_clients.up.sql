-- 1. Add the column as nullable first (without a default)
ALTER TABLE clients ADD COLUMN email TEXT;

-- 2. Give existing clients a unique temporary email using their UUID
-- This prevents the "duplicate key" error
UPDATE clients SET email = 'client_' || id || '@example.com' WHERE email IS NULL;

-- 3. Now make it NOT NULL
ALTER TABLE clients ALTER COLUMN email SET NOT NULL;

-- 4. Now we can safely add the unique constraint
ALTER TABLE clients ADD CONSTRAINT unique_trainer_client_email UNIQUE (trainer_id, email);