package entities

import "errors"

var (
	ErrPRMerged              = errors.New("pull request is already merged")
	ErrReviewerNotAssigned   = errors.New("reviewer is not assigned to this PR")
	ErrUserNoID           = errors.New("user: no user_id")
	ErrUserNoUsername     = errors.New("user: no username")
	ErrTeamNoName         = errors.New("team: no team name")
	ErrTeamPresent        = errors.New("team: user already in this team")
	ErrPRNoID             = errors.New("pr: no pull_request_id")
	ErrPRNoName           = errors.New("pr: no pull_request_name")
	ErrPRNoAuthor         = errors.New("pr: no author_id")
	ErrPRTooManyReviewers = errors.New("pr: too many reviewers")
	ErrAuthorIsReviewer   = errors.New("pr: can't assign author as reviewer")
)
