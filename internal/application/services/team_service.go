package services

import (
	"context"
	"errors"

	"github.com/Traunin/review-assigner/internal/application/dto"
	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/domain/repositories"
)

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
)

type TeamService interface {
    CreateTeam(ctx context.Context, cmd dto.CreateTeamCmd) (*dto.TeamDTO, error)
    GetTeam(ctx context.Context, teamName string) (*dto.TeamDTO, error)
}

type teamService struct {
	teams repositories.TeamRepository
	users repositories.UserRepository
}

func NewTeamService(
	teams repositories.TeamRepository,
	users repositories.UserRepository,
) TeamService {
	return &teamService{
		teams: teams,
		users: users,
	}
}

func (s *teamService) CreateTeam(
	ctx context.Context,
	cmd dto.CreateTeamCmd,
) (*dto.TeamDTO, error) {
	existing, err := s.teams.FindByName(ctx, cmd.TeamName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTeamExists
	}

	team, err := entities.NewTeam(cmd.TeamName, 0)
	if err != nil {
		return nil, err
	}

	if err := s.teams.Create(ctx, team); err != nil {
		return nil, err
	}

	team, err = s.teams.FindByName(ctx, cmd.TeamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("failed to load created team")
	}

	for _, m := range cmd.Members {
		u, err := s.users.FindByID(ctx, entities.UserID(m.UserID))
		if err != nil {
			return nil, err
		}

		tid := team.ID()
		if u == nil {
			newUser, err := entities.NewUser(
				entities.UserID(m.UserID),
				m.Username,
				m.IsActive,
				&tid,
			)
			if err != nil {
				return nil, err
			}

			if err := s.users.Create(ctx, newUser); err != nil {
				return nil, err
			}

		} else {
			u.SetUsername(m.Username)
			u.SetActive(m.IsActive)
			u.SetTeamID(&tid)

			if err := s.users.Update(ctx, u); err != nil {
				return nil, err
			}
		}
	}

	return &dto.TeamDTO{
		ID:       int64(team.ID()),
		TeamName: team.Name(),
	}, nil
}

func (s *teamService) GetTeam(
	ctx context.Context,
	teamName string,
) (*dto.TeamDTO, error) {

	team, err := s.teams.FindByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	return &dto.TeamDTO{
		ID:       int64(team.ID()),
		TeamName: team.Name(),
	}, nil
}
