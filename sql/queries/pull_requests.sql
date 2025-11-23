-- name: CreatePullRequest :one
INSERT INTO pull_requests (
    pull_request_id, 
    pull_request_name, 
    author_id, 
    status,
    created_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at;

-- name: GetPullRequestByID :one
SELECT 
    pull_request_id, 
    pull_request_name, 
    author_id, 
    status,
    created_at,
    merged_at
FROM pull_requests
WHERE pull_request_id = $1;

-- name: UpdatePRStatus :one
UPDATE pull_requests
SET 
    status = $2,
    merged_at = $3
WHERE pull_request_id = $1
RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at;

-- name: PRExists :one
SELECT EXISTS(
    SELECT 1 FROM pull_requests WHERE pull_request_id = $1
);

-- name: GetPRsByAuthor :many
SELECT 
    pull_request_id, 
    pull_request_name, 
    author_id, 
    status,
    created_at,
    merged_at
FROM pull_requests
WHERE author_id = $1
ORDER BY created_at DESC;

-- name: GetOpenPRs :many
SELECT 
    pull_request_id, 
    pull_request_name, 
    author_id, 
    status,
    created_at,
    merged_at
FROM pull_requests
WHERE status = 'OPEN'
ORDER BY created_at DESC;

-- name: DeletePullRequest :exec
DELETE FROM pull_requests
WHERE pull_request_id = $1;

-- name: GetPullRequests :many
SELECT 
    pull_request_id, 
    pull_request_name, 
    author_id, 
    status,
    created_at,
    merged_at
FROM pull_requests
ORDER BY created_at DESC;
