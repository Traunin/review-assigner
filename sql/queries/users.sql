-- name: CreateUser :one
INSERT INTO users (user_id, username, is_active, team_id)
VALUES ($1, $2, $3, $4)
RETURNING user_id, username, is_active, team_id;

-- name: GetUserByID :one
SELECT user_id, username, is_active, team_id
FROM users
WHERE user_id = $1;

-- name: GetUsers :many
SELECT user_id, username, is_active, team_id
FROM users;

-- name: UpdateUserStatus :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING user_id, username, is_active;

-- name: GetUsersByTeamID :many
SELECT user_id, username, is_active, team_id
FROM users
WHERE team_id = $1;

-- name: GetTeamByUserID :one
SELECT t.id, t.team_name FROM teams t
JOIN users u ON u.team_id = t.id
WHERE u.user_id = $1;

-- name: GetActiveUsersByTeamID :many
SELECT user_id, username, is_active, team_id
FROM users
WHERE team_id = $1 AND is_active = true;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE user_id = $1
);

-- name: GetActiveUsers :many
SELECT user_id, username, is_active, team_id
FROM users
WHERE is_active = true;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, is_active = $3, team_id = $4
WHERE user_id = $1;
