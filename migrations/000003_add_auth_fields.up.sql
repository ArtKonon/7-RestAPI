ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS password_hash TEXT;

UPDATE users
SET role = 'reader'
WHERE role IS NULL OR role = 'user';

ALTER TABLE IF EXISTS users
    ALTER COLUMN role SET DEFAULT 'reader';

ALTER TABLE IF EXISTS users
    DROP CONSTRAINT IF EXISTS chk_users_role;

ALTER TABLE IF EXISTS users
    ADD CONSTRAINT chk_users_role
    CHECK (role IN ('reader', 'librarian', 'admin'));

UPDATE users
SET password_hash = ''
WHERE password_hash IS NULL;
