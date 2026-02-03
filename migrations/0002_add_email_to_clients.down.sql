ALTER TABLE clients DROP CONSTRAINT IF EXISTS unique_trainer_client_email;
ALTER TABLE clients DROP COLUMN IF EXISTS email;