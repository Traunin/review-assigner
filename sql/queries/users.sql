-- name: CreateUser :one
INSERT INTO users (user_id, username, is_active)
VALUES ($1, $2, $3)
RETURNING user_id, username, is_active;

-- name: GetUserByID :one
SELECT user_id, username, is_active
FROM users
WHERE user_id = $1;

-- name: GetUsers :many
SELECT user_id, username, is_active
FROM users;

-- name: UpdateUserStatus :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING user_id, username, is_active;

-- name: GetActiveTeamMembers :many
SELECT u.user_id, u.username, u.is_active
FROM users u
JOIN team_members tm ON u.user_id = tm.user_id
WHERE tm.team_id = $1 AND u.is_active = true;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE user_id = $1
);

-- name: GetActiveUsers :many
SELECT user_id, username, is_active
FROM users
WHERE is_active = true;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, is_active = $3
WHERE user_id = $1;
