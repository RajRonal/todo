CREATE TYPE status AS ENUM ('draft', 'active');
ALTER TABLE todo
    DROP COLUMN status;
ALTER TABLE todo
    ADD    is_complete BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE todo
    ADD    archived_at     TIMESTAMP with time zone;
ALTER TABLE todo
    ADD    status     status;