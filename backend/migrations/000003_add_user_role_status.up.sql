ALTER TABLE users
ADD COLUMN role VARCHAR(30) NOT NULL DEFAULT 'user',
ADD COLUMN status VARCHAR(30) NOT NULL DEFAULT 'active';

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);