-- Migration: 005_seed_default_stages_function
-- Creates a helper function that inserts the default pipeline stages for a job.
-- Called automatically via trigger when a new job is inserted.

CREATE OR REPLACE FUNCTION insert_default_job_stages()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO job_stages (job_id, name, order_index) VALUES
        (NEW.id, 'Screening',           1),
        (NEW.id, 'Technical Challenge', 2),
        (NEW.id, 'Team Interview',       3),
        (NEW.id, 'Manager Interview',    4),
        (NEW.id, 'Offer',               5),
        (NEW.id, 'Hired',               6);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_default_job_stages
AFTER INSERT ON jobs
FOR EACH ROW
EXECUTE FUNCTION insert_default_job_stages();
