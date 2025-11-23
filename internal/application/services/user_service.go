package services

import (
	"context"
	"errors"

	"github.com/Traunin/review-assigner/internal/application/dto"
	"github.com/Traunin/review-assigner/internal/application/mapper"
	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/domain/repositories"
)

type UserService interface {
	AddUser(ctx context.Context, user *dto.UserDTO) error
	DeleteUserByID(ctx context.Context, userID entities.UserID) error
	GetAll(ctx context.Context) ([]*dto.UserDTO, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepository,
	}
}

func (s *userService) AddUser(
	ctx context.Context,
	user *dto.UserDTO,
) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	var teamID *entities.TeamID

	if user.TeamID != nil {
		intValue := *user.TeamID
		convertedTeamID := entities.TeamID(intValue)
		teamID = &convertedTeamID
	}

	entity, err := entities.NewUser(
		user.UserID,
		user.Username,
		user.IsActive,
		teamID,
	)
	if err != nil {
		return err
	}
	err = s.userRepo.Create(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUserByID(
	ctx context.Context,
	id entities.UserID,
) error {
	err := s.userRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetAll(ctx context.Context) ([]*dto.UserDTO, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToUserDTOs(users), nil
}
