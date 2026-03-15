CREATE TABLE IF NOT EXISTS permissions (
    id         VARCHAR(36) PRIMARY KEY,
    file_id    VARCHAR(36) NOT NULL REFERENCES files (id) ON DELETE CASCADE,
    user_id    VARCHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    can_read   BOOLEAN     NOT NULL DEFAULT FALSE,
    can_write  BOOLEAN     NOT NULL DEFAULT FALSE,
    can_delete BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (file_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_permissions_file_id ON permissions (file_id);
CREATE INDEX IF NOT EXISTS idx_permissions_user_id ON permissions (user_id);