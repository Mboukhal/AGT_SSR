-- name: CreateUser :exec
INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id;

-- name: GetUserByEmail :one
SELECT id, email, password, username, is_active, created_at, updated_at FROM users WHERE email = $1;

-- name: CheckUserEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);