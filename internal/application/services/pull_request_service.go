package services

import (
	"context"
	"errors"
	"time"

	"github.com/Traunin/review-assigner/internal/application/dto"
	"github.com/Traunin/review-assigner/internal/application/mapper"
	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/domain/repositories"
	ds "github.com/Traunin/review-assigner/internal/domain/services"
)

var (
	ErrNotFound = errors.New("pull request not found")
	ErrPRNotLoaded = errors.New("failed to load pull request")
)

type PullRequestService interface {
    GetByID(ctx context.Context, id entities.PullRequestID) (dto.PullRequestDTO, error)
    Create(ctx context.Context, input dto.CreatePRCmd) (dto.PullRequestDTO, error)
	Merge(ctx context.Context, id entities.PullRequestID) (dto.PullRequestDTO, error)
	ReassignReviewer(ctx context.Context, input dto.ReassignReviewerCmd) (*dto.ReassignedDTO, error)
}

type pullRequestService struct {
    repo repositories.PullRequestRepository
	prService ds.ReviewerAssignmentService
}

func NewPullRequestService(
	repo repositories.PullRequestRepository,
	prService ds.ReviewerAssignmentService,
) PullRequestService {
    return &pullRequestService{
		repo: repo,
		prService: prService,
	}
}

func (s *pullRequestService) GetByID(
	ctx context.Context,
	id entities.PullRequestID,
) (dto.PullRequestDTO, error) {
    pr, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return dto.PullRequestDTO{}, err
    }
    if pr == nil {
        return dto.PullRequestDTO{}, ErrNotFound
    }

    return mapper.ToPullRequestDTO(pr), nil
}

func (s *pullRequestService) Create(
	ctx context.Context,
	input dto.CreatePRCmd,
) (dto.PullRequestDTO, error) {
	entity, err := entities.NewPullRequest(
		input.PullRequestID,
		input.PullRequestName,
		input.AuthorID,
		entities.StatusOpen,
		nil,
		time.Now(),
		nil,
	)
	
	if err != nil {
		return dto.PullRequestDTO{},  err
	}
	created, err := s.prService.CreateAndAssign(ctx, entity)

	if err != nil {
		return dto.PullRequestDTO{}, err
	}
	if created == nil {
		return dto.PullRequestDTO{}, ErrPRNotLoaded
	}

	return mapper.ToPullRequestDTO(created), nil
}

func (s *pullRequestService) Merge(
	ctx context.Context,
	id entities.PullRequestID,
) (dto.PullRequestDTO, error) {
	pr, err := s.prService.Merge(ctx, id)
	if err != nil {
		return dto.PullRequestDTO{}, err
	}

	return mapper.ToPullRequestDTO(pr), nil
}

func (s *pullRequestService) ReassignReviewer(
	ctx context.Context,
	input dto.ReassignReviewerCmd,
) (*dto.ReassignedDTO, error) {
	assigned, pr, err := s.prService.ReassignReviewer(
		ctx,
		input.PullRequestID,
		input.OldUserID,
	)
	if err != nil {
		return nil, err
	}

	prDTO := mapper.ToPullRequestDTO(pr)
	return &dto.ReassignedDTO{PullRequestID: &prDTO, Assigned: assigned}, nil
}
