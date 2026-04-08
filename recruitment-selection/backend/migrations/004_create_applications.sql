-- Migration: 004_create_applications
-- Tracks a candidate's application to a job, including the current pipeline
-- stage and the overall decision status.
-- A candidate can only apply once per job (uq_application).

CREATE TYPE application_status AS ENUM (
    'pending',
    'in_progress',
    'accepted',
    'rejected',
    'withdrawn'
);

CREATE TABLE IF NOT EXISTS applications (
    id               UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id           UUID               NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    candidate_id     UUID               NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_stage_id UUID               REFERENCES job_stages(id) ON DELETE SET NULL,
    status           application_status NOT NULL DEFAULT 'pending',
    cover_letter     TEXT,
    cv_url           VARCHAR(500),
    created_at       TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_application UNIQUE (job_id, candidate_id)
);

CREATE INDEX idx_applications_job_id       ON applications(job_id);
CREATE INDEX idx_applications_candidate_id ON applications(candidate_id);
CREATE INDEX idx_applications_status       ON applications(status);
