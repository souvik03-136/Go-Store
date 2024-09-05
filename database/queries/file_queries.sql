-- name: CreateFile :exec
INSERT INTO files (id, name, path, size, content_type, owner_id, created_at, updated_at)
VALUES (@id, @name, @path, @size, @content_type, @owner_id, NOW(), NOW());

-- name: GetFileByID :one
SELECT * FROM files WHERE id = @id;

-- name: UpdateFile :exec
UPDATE files
SET name = @name, path = @path, size = @size, content_type = @content_type, updated_at = NOW()
WHERE id = @id;

-- name: DeleteFile :exec
DELETE FROM files WHERE id = @id;

-- name: ListFilesByOwner :many
SELECT * FROM files WHERE owner_id = @owner_id;
