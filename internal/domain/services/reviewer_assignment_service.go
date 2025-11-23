package services

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/domain/repositories"
)

var (
	ErrPRAlreadyMerged = errors.New("pull request already merged")
	ErrTeamNotFound    = errors.New("team not found")
	ErrPRNotFound      = errors.New("pull request not found")
	ErrUserNotReviewer = errors.New("user is not a reviewer")
	ErrPRAlreadyExists = errors.New("pull request already exists")
	ErrAuthorNotFound  = errors.New("author not found")
	ErrNoCandidate     = errors.New("no candidate")
)

type ReviewerAssignmentService interface {
	CreateAndAssign(
		ctx context.Context,
		pr *entities.PullRequest,
	) (*entities.PullRequest, error)
	ReassignReviewer(
		ctx context.Context,
		prID entities.PullRequestID,
		oldReviewerID entities.UserID,
	) (entities.UserID, *entities.PullRequest, error)
	Merge(
		ctx context.Context,
		prID entities.PullRequestID,
	) (*entities.PullRequest, error)
}

type reviewerAssignmentService struct {
	prRepo   repositories.PullRequestRepository
	userRepo repositories.UserRepository
	teamRepo repositories.TeamRepository
}

func NewReviewerAssignmentService(
	userRepo repositories.UserRepository,
	prRepo repositories.PullRequestRepository,
	teamRepo repositories.TeamRepository,
) ReviewerAssignmentService {
	return &reviewerAssignmentService{
		userRepo: userRepo,
		prRepo:   prRepo,
		teamRepo: teamRepo,
	}
}

func (s *reviewerAssignmentService) CreateAndAssign(
	ctx context.Context,
	pr *entities.PullRequest,
) (*entities.PullRequest, error) {
	existing, err := s.prRepo.FindByID(ctx, pr.ID())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPRAlreadyExists
	}

	author, err := s.userRepo.FindByID(ctx, pr.AuthorID())
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, ErrAuthorNotFound
	}

	team, err := s.teamRepo.FindByUserID(ctx, author.ID())
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	activeMembers, err := s.teamRepo.FindActiveReviewersByTeamID(ctx, team.ID())
	if err != nil {
		return nil, err
	}

	reviewerIDs := selectReviewers(
		activeMembers,
		author.ID(),
		nil,
		entities.MaxReviewers,
	)

	for _, rid := range reviewerIDs {
		err := pr.AssignReviewer(rid)
		if err != nil {
			return nil, err
		}
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func selectReviewers(
	activeMembers []*entities.User,
	authorID entities.UserID,
	excludeUserIDs []entities.UserID,
	maxReviewers int,
) []entities.UserID {
	excluded := make(map[entities.UserID]bool)
	excluded[authorID] = true
	for _, uid := range excludeUserIDs {
		excluded[uid] = true
	}

	candidates := make([]entities.UserID, 0)
	for _, member := range activeMembers {
		if !excluded[member.ID()] {
			candidates = append(candidates, member.ID())
		}
	}

	shuffle(candidates)
	selectCount := min(len(candidates), maxReviewers)
	return candidates[:selectCount]
}

func selectReplacementReviewer(
	activeMembers []*entities.User,
	authorID entities.UserID,
	excludeUserIDs []entities.UserID,
) (entities.UserID, error) {
	selected := selectReviewers(activeMembers, authorID, excludeUserIDs, 1)
	if len(selected) == 0 {
		return "", ErrNoCandidate
	}
	return selected[0], nil
}

func (s *reviewerAssignmentService) ReassignReviewer(
	ctx context.Context,
	prID entities.PullRequestID,
	oldReviewerID entities.UserID,
) (entities.UserID, *entities.PullRequest, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return "", nil, err
	}
	if pr == nil {
		return "", nil, ErrPRNotFound
	}
	if pr.IsMerged() {
		return "", nil, ErrPRAlreadyMerged
	}

	if !pr.HasReviewer(oldReviewerID) {
		return "", nil, ErrUserNotReviewer
	}

	team, err := s.teamRepo.FindByUserID(ctx, pr.AuthorID())
	if err != nil {
		return "", nil, err
	}
	if team == nil {
		return "", nil, ErrTeamNotFound
	}

	active, err := s.teamRepo.FindActiveReviewersByTeamID(ctx, team.ID())
	if err != nil {
		return "", nil, err
	}

	exclude := make([]entities.UserID, 0, len(pr.ReviewerIDs())+2)
	exclude = append(exclude, pr.ReviewerIDs()...)
	exclude = append(exclude, pr.AuthorID(), oldReviewerID)

	newReviewerID, err := selectReplacementReviewer(
		active,
		pr.AuthorID(),
		exclude,
	)
	if err != nil {
		return "", nil, err
	}

	err = pr.UnassignReviewer(oldReviewerID)
	if err != nil {
		return "", nil, err
	}
	err = pr.AssignReviewer(newReviewerID)
	if err != nil {
		return "", nil, err
	}

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return "", nil, err
	}

	return newReviewerID, pr, nil
}

func (s *reviewerAssignmentService) Merge(
	ctx context.Context,
	prID entities.PullRequestID,
) (*entities.PullRequest, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.IsMerged() {
		return pr, nil
	}

	pr.Merge()

	if err := s.prRepo.Update(ctx, pr); err != nil {
		// race: someone else merged it first
		latest, findErr := s.prRepo.FindByID(ctx, prID)
		if findErr == nil && latest != nil && latest.IsMerged() {
			return latest, nil
		}
		return nil, err
	}

	return pr, nil
}

// Fisher-Yates shuffle to avoid rand seed shenanigans
func shuffle(slice []entities.UserID) {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		// as per docs, "Int cannot return an error when using rand.Reader"
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
