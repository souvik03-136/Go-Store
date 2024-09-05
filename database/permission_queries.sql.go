// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: permission_queries.sql

package database

import (
	"context"
)

const createPermission = `-- name: CreatePermission :exec
INSERT INTO permissions (id, file_id, user_id, can_read, can_write, can_delete, created_at, updated_at)
VALUES (@id, @file_id, @user_id, @can_read, @can_write, @can_delete, NOW(), NOW())
`

func (q *Queries) CreatePermission(ctx context.Context) error {
	_, err := q.db.Exec(ctx, createPermission)
	return err
}

const deletePermission = `-- name: DeletePermission :exec
DELETE FROM permissions WHERE file_id = @file_id AND user_id = @user_id
`

func (q *Queries) DeletePermission(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deletePermission)
	return err
}

const getPermissionsByFileID = `-- name: GetPermissionsByFileID :many
SELECT id, file_id, user_id, can_read, can_write, can_delete, created_at, updated_at FROM permissions WHERE file_id = @file_id
`

func (q *Queries) GetPermissionsByFileID(ctx context.Context) ([]Permission, error) {
	rows, err := q.db.Query(ctx, getPermissionsByFileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Permission
	for rows.Next() {
		var i Permission
		if err := rows.Scan(
			&i.ID,
			&i.FileID,
			&i.UserID,
			&i.CanRead,
			&i.CanWrite,
			&i.CanDelete,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePermission = `-- name: UpdatePermission :exec
UPDATE permissions
SET can_read = @can_read, can_write = @can_write, can_delete = @can_delete, updated_at = NOW()
WHERE file_id = @file_id AND user_id = @user_id
`

func (q *Queries) UpdatePermission(ctx context.Context) error {
	_, err := q.db.Exec(ctx, updatePermission)
	return err
}
