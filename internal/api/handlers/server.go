package handlers

import (
	"github.com/Traunin/review-assigner/internal/application/services"
	"github.com/Traunin/review-assigner/internal/domain/repositories"
)

type Server struct {
	teamService services.TeamService
	prService   services.PullRequestService
	userRepo    repositories.UserRepository
	teamRepo    repositories.TeamRepository
	prRepo      repositories.PullRequestRepository
}

func NewServer(
	teamService services.TeamService,
	prService services.PullRequestService,
	userRepo repositories.UserRepository,
	teamRepo repositories.TeamRepository,
	prRepo repositories.PullRequestRepository,
) *Server {
	return &Server{
		teamService: teamService,
		prService:   prService,
		userRepo:    userRepo,
		teamRepo:    teamRepo,
		prRepo:      prRepo,
	}
}
