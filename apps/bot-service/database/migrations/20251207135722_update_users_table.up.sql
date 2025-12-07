-- Rename label column to name
ALTER TABLE users RENAME COLUMN label TO name;

-- Remove columns
ALTER TABLE users
    DROP COLUMN assigned_to,
    DROP COLUMN notes;

-- Add new columns for user details
ALTER TABLE users
    ADD COLUMN job_title VARCHAR(255),
    ADD COLUMN gender VARCHAR(20),
    ADD COLUMN date_of_birth DATE;
