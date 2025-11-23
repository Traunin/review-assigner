package entities

import (
	"slices"
	"time"
)

const MaxReviewers = 2

type Reviewer struct {
	UserID     UserID
	AssignedAt time.Time
}

type PullRequest struct {
	id        PullRequestID
	name      string
	authorID  UserID
	status    PRStatus
	reviewers []Reviewer
	createdAt time.Time
	mergedAt  *time.Time
}

func NewPullRequest(
	id PullRequestID,
	name string,
	authorID UserID,
	status PRStatus,
	reviewers []Reviewer,
	createdAt time.Time,
	mergedAt *time.Time,
) (*PullRequest, error) {
	for _, reviewer := range reviewers {
		if reviewer.UserID == authorID {
			return nil, ErrAuthorIsReviewer
		}
	}

	pr := &PullRequest{
		id:        id,
		name:      name,
		authorID:  authorID,
		status:    status,
		createdAt: createdAt,
		mergedAt:  mergedAt,
		reviewers: reviewers,
	}

	if err := pr.validate(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (pr *PullRequest) validate() error {
	if pr.id == "" {
		return ErrPRNoID
	}
	if pr.name == "" {
		return ErrPRNoName
	}
	if pr.authorID == "" {
		return ErrPRNoAuthor
	}
	return nil
}

func (pr *PullRequest) Merge() {
	if pr.status == StatusMerged {
		return
	}
	pr.status = StatusMerged
	now := time.Now()
	pr.mergedAt = &now
}

func (pr *PullRequest) HasEnoughReviewers() bool {
	return len(pr.reviewers) == MaxReviewers
}

func (pr *PullRequest) IsAssigneeValid(id UserID) error {
	if pr.HasEnoughReviewers() {
		return ErrPRTooManyReviewers
	}

	if pr.status == StatusMerged {
		return ErrPRMerged
	}

	if pr.authorID == id {
		return ErrAuthorIsReviewer
	}

	return nil
}

func (pr *PullRequest) AssignReviewer(id UserID) error {
	if err := pr.IsAssigneeValid(id); err != nil {
		return err
	}

	if slices.Contains(pr.reviewers, Reviewer{UserID: id}) {
		return ErrTeamPresent
	}

	pr.reviewers = append(pr.reviewers, Reviewer{
		UserID:     id,
		AssignedAt: time.Now(),
	})

	return nil
}

func (pr *PullRequest) ReassignReviewer(
	oldUserID UserID,
	newUserID UserID,
) error {
	if pr.status == StatusMerged {
		return ErrPRMerged
	}

	idx := -1
	for i, r := range pr.reviewers {
		if r.UserID == oldUserID {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrReviewerNotAssigned
	}

	pr.reviewers[idx] = Reviewer{
		UserID:     newUserID,
		AssignedAt: time.Now(),
	}

	return nil
}

func (pr *PullRequest) UnassignReviewer(id UserID) error {
	if pr.status == StatusMerged {
		return ErrPRMerged
	}

	pr.reviewers = slices.DeleteFunc(pr.reviewers, func(r Reviewer) bool {
		return r.UserID == id
	})

	return nil
}

func (pr *PullRequest) ReviewerIDs() []UserID {
	ids := make([]UserID, len(pr.reviewers))
	for i, r := range pr.reviewers {
		ids[i] = r.UserID
	}
	return ids
}

func (pr *PullRequest) ID() PullRequestID {
	return pr.id
}

func (pr *PullRequest) Name() string {
	return pr.name
}

func (pr *PullRequest) AuthorID() UserID {
	return pr.authorID
}

func (pr *PullRequest) Status() PRStatus {
	return pr.status
}

func (pr *PullRequest) CreatedAt() time.Time {
	return pr.createdAt
}

func (pr *PullRequest) MergedAt() time.Time {
	if pr.mergedAt == nil {
		return time.Time{}
	}
	return *pr.mergedAt
}

func (pr *PullRequest) MergedAtPtr() *time.Time {
    return pr.mergedAt
}

func (pr *PullRequest) Reviewers() []Reviewer {
	return slices.Clone(pr.reviewers)
}
