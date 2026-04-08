-- Migration: 003_create_job_stages
-- Stores the pipeline stages defined by the recruiter for each job.
-- order_index must be unique per job to maintain a clear sequence.

CREATE TABLE IF NOT EXISTS job_stages (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id      UUID        NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    order_index INT         NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_job_stage_order UNIQUE (job_id, order_index)
);

CREATE INDEX idx_job_stages_job_id ON job_stages(job_id);
