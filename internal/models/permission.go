package models

import "time"

// Permission represents a user's permissions on a file.
type Permission struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`    // References the file the permission applies to
	UserID    string    `json:"user_id"`    // References the user who has the permission
	CanRead   bool      `json:"can_read"`   // Whether the user can read the file
	CanWrite  bool      `json:"can_write"`  // Whether the user can write to the file
	CanDelete bool      `json:"can_delete"` // Whether the user can delete the file
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPermission creates a new Permission instance.
func NewPermission(id, fileID, userID string, canRead, canWrite, canDelete bool) *Permission {
	return &Permission{
		ID:        id,
		FileID:    fileID,
		UserID:    userID,
		CanRead:   canRead,
		CanWrite:  canWrite,
		CanDelete: canDelete,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdatePermission updates the permission details.
func (p *Permission) UpdatePermission(canRead, canWrite, canDelete bool) error {
	p.CanRead = canRead
	p.CanWrite = canWrite
	p.CanDelete = canDelete
	p.UpdatedAt = time.Now()
	return nil
}

// CanAccess checks if the user has the required permissions to access the file.
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

// GrantFullAccess gives full (read, write, delete) access to a user.
func (p *Permission) GrantFullAccess() {
	p.CanRead = true
	p.CanWrite = true
	p.CanDelete = true
	p.UpdatedAt = time.Now()
}

// RevokeAllAccess revokes all permissions for a user.
func (p *Permission) RevokeAllAccess() {
	p.CanRead = false
	p.CanWrite = false
	p.CanDelete = false
	p.UpdatedAt = time.Now()
}
