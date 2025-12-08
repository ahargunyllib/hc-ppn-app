-- Remove new columns
ALTER TABLE users
    DROP COLUMN date_of_birth,
    DROP COLUMN gender,
    DROP COLUMN job_title;

-- Add back removed columns
ALTER TABLE users
    ADD COLUMN assigned_to VARCHAR(255),
    ADD COLUMN notes TEXT;

-- Rename name column back to label
ALTER TABLE users RENAME COLUMN name TO label;
