-- name: AddReviewer :exec
INSERT INTO reviewers (pull_request_id, user_id, assigned_at)
VALUES ($1, $2, $3);

-- name: RemoveReviewer :exec
DELETE FROM reviewers
WHERE pull_request_id = $1 AND user_id = $2;

-- name: GetReviewersByPR :many
SELECT user_id, assigned_at
FROM reviewers
WHERE pull_request_id = $1
ORDER BY assigned_at;

-- name: GetPRsByReviewer :many
SELECT 
    pr.pull_request_id,
    pr.pull_request_name,
    pr.author_id,
    pr.status
FROM pull_requests pr
JOIN reviewers rev ON pr.pull_request_id = rev.pull_request_id
WHERE rev.user_id = $1
ORDER BY pr.created_at DESC;

-- name: IsUserReviewer :one
SELECT EXISTS(
    SELECT 1 
    FROM reviewers 
    WHERE pull_request_id = $1 AND user_id = $2
);

-- name: GetReviewerCount :one
SELECT COUNT(*)
FROM reviewers
WHERE pull_request_id = $1;

-- name: ReplaceReviewer :exec
WITH deleted AS (
    DELETE FROM reviewers r
    WHERE r.pull_request_id = $1 AND r.user_id = $2
)
INSERT INTO reviewers (pull_request_id, user_id, assigned_at)
VALUES ($1, $3, $4)
ON CONFLICT DO NOTHING;
