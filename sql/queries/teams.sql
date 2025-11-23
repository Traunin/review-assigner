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

-- name: GetTeamWithMembers :many
SELECT 
    t.id as team_id,
    t.team_name,
    u.user_id,
    u.username,
    u.is_active
FROM teams t
LEFT JOIN team_members tm ON t.id = tm.team_id
LEFT JOIN users u ON tm.user_id = u.user_id
WHERE t.team_name = $1;

-- name: AddTeamMember :exec
INSERT INTO team_members (team_id, user_id)
VALUES ($1, $2)
ON CONFLICT (team_id, user_id) DO NOTHING;

-- name: RemoveTeamMember :exec
DELETE FROM team_members
WHERE team_id = $1 AND user_id = $2;

-- name: GetTeamMemberCount :one
SELECT COUNT(*) 
FROM team_members 
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

-- name: DeleteTeamMembers :exec
DELETE FROM team_members
WHERE team_id = $1;

-- name: GetTeamsByUserID :many
SELECT t.id, t.team_name
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = $1;
