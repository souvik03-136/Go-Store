-- name: CreatePermission :exec
INSERT INTO permissions (id, file_id, user_id, can_read, can_write, can_delete, created_at, updated_at)
VALUES (@id, @file_id, @user_id, @can_read, @can_write, @can_delete, NOW(), NOW());

-- name: GetPermissionsByFileID :many
SELECT * FROM permissions WHERE file_id = @file_id;

-- name: UpdatePermission :exec
UPDATE permissions
SET can_read = @can_read, can_write = @can_write, can_delete = @can_delete, updated_at = NOW()
WHERE file_id = @file_id AND user_id = @user_id;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE file_id = @file_id AND user_id = @user_id;
