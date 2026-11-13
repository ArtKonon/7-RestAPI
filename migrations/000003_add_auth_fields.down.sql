ALTER TABLE IF EXISTS users
    DROP CONSTRAINT IF EXISTS chk_users_role;

UPDATE users
SET role = 'user'
WHERE role NOT IN ('user', 'admin');

ALTER TABLE IF EXISTS users
    ALTER COLUMN role SET DEFAULT 'user';

ALTER TABLE IF EXISTS users
    ADD CONSTRAINT chk_users_role
    CHECK (role IN ('user', 'admin'));

ALTER TABLE IF EXISTS users
    DROP COLUMN IF EXISTS password_hash;
