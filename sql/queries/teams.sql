-- name: CreateTeam :one
INSERT INTO teams (team_name)
VALUES ($1)
RETURNING id, team_name;

-- name: GetTeamByName :one
SELECT id, team_name
FROM teams
WHERE team_name = $1;

-- name: GetTeamByID :one
SELECT id, team_name
FROM teams
WHERE id = $1;

-- name: TeamExists :one
SELECT EXISTS(
    SELECT 1 FROM teams WHERE team_name = $1
);

-- name: GetTeamMemberCount :one
SELECT COUNT(*) 
FROM users 
WHERE team_id = $1;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;

-- name: GetTeams :many
SELECT id, team_name
FROM teams;

-- name: UpdateTeam :exec
UPDATE teams
SET team_name = $2
WHERE id = $1;
