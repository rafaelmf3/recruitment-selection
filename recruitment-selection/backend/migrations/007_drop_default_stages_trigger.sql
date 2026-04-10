-- Migration: 007_drop_default_stages_trigger
-- The trigger created in 005 inserted English default stages automatically on
-- every job insert, which conflicts with explicit stage management done by the
-- application (CreateJob / UpdateJobStages) and the seed script.
-- Dropping it so stages are only created when explicitly provided.

DROP TRIGGER IF EXISTS trg_default_job_stages ON jobs;
DROP FUNCTION IF EXISTS insert_default_job_stages();
