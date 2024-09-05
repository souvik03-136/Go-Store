-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES (@id, @username, @email, @password_hash, NOW(), NOW());

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = @username;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email;

-- name: UpdateUser :exec
UPDATE users
SET username = @username, email = @email, password_hash = @password_hash, updated_at = NOW()
WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;
