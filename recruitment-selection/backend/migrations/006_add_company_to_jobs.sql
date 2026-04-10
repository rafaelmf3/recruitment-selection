-- Migration: 006_add_company_to_jobs
-- Adds the company column to the jobs table so recruiters can associate
-- a job posting with the company name they are hiring for.

ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS company VARCHAR(255);
