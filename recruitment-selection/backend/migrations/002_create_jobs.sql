-- Migration: 002_create_jobs
-- Creates the job_status enum type and jobs table.
--
-- Status state machine:
--   open  --> paused | closed | cancelled
--   paused --> open  | closed | cancelled
--   closed    (terminal - process completed, someone was hired)
--   cancelled (terminal - process abandoned, no hire)

CREATE TYPE job_status AS ENUM ('open', 'paused', 'closed', 'cancelled');

CREATE TABLE IF NOT EXISTS jobs (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    recruiter_id UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    description  TEXT         NOT NULL,
    requirements TEXT,
    location     VARCHAR(255),
    salary_min   NUMERIC(12, 2),
    salary_max   NUMERIC(12, 2),
    status       job_status   NOT NULL DEFAULT 'open',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_salary_range CHECK (
        salary_min IS NULL OR salary_max IS NULL OR salary_min <= salary_max
    )
);

CREATE INDEX idx_jobs_recruiter_id ON jobs(recruiter_id);
CREATE INDEX idx_jobs_status       ON jobs(status);
CREATE INDEX idx_jobs_salary       ON jobs(salary_min, salary_max);
