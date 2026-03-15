
CREATE TABLE IF NOT EXISTS files (
    id           VARCHAR(36)   PRIMARY KEY,
    name         VARCHAR(255)  NOT NULL,
    path         VARCHAR(1024) NOT NULL,
    url          VARCHAR(2048) NOT NULL,
    size         BIGINT        NOT NULL DEFAULT 0,
    content_type VARCHAR(255)  NOT NULL DEFAULT '',
    owner_id     VARCHAR(36)   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_files_owner_id ON files (owner_id);