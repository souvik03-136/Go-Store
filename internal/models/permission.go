// internal/models/permission.go

package models

import "time"

// Permission represents a user's access rights on a specific file.
type Permission struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	UserID    string    `json:"user_id"`
	CanRead   bool      `json:"can_read"`
	CanWrite  bool      `json:"can_write"`
	CanDelete bool      `json:"can_delete"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPermission constructs a Permission with the current timestamp.
func NewPermission(id, fileID, userID string, canRead, canWrite, canDelete bool) *Permission {
	now := time.Now()
	return &Permission{
		ID:        id,
		FileID:    fileID,
		UserID:    userID,
		CanRead:   canRead,
		CanWrite:  canWrite,
		CanDelete: canDelete,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CanAccess checks if the permission grants the requested access type.
// Valid values for permissionType: "read", "write", "delete".
func (p *Permission) CanAccess(permissionType string) bool {
	switch permissionType {
	case "read":
		return p.CanRead
	case "write":
		return p.CanWrite
	case "delete":
		return p.CanDelete
	default:
		return false
	}
}

// GrantFullAccess enables all permission flags.
func (p *Permission) GrantFullAccess() {
	p.CanRead = true
	p.CanWrite = true
	p.CanDelete = true
	p.UpdatedAt = time.Now()
}

// RevokeAllAccess disables all permission flags.
func (p *Permission) RevokeAllAccess() {
	p.CanRead = false
	p.CanWrite = false
	p.CanDelete = false
	p.UpdatedAt = time.Now()
}
