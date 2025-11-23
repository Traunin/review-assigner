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
	id        PrID
	name      string
	authorID  UserID
	status    PRStatus
	reviewers []Reviewer
	createdAt time.Time
	mergedAt  *time.Time
}

func NewPullRequest(
	id PrID,
	name string,
	authorID UserID,
	assignedUserIDs []UserID,
) (*PullRequest, error) {
	pr := &PullRequest{
		id:        id,
		name:      name,
		authorID:  authorID,
		status:    StatusOpen,
		createdAt: time.Now(),
		reviewers: make([]Reviewer, 0, len(assignedUserIDs)),
	}

	if err := pr.validate(); err != nil {
		return nil, err
	}

	for _, uid := range assignedUserIDs {
		if uid == authorID {
			return nil, ErrAuthorIsReviewer
		}
		pr.reviewers = append(pr.reviewers, Reviewer{
			UserID:     uid,
			AssignedAt: time.Now(),
		})
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
