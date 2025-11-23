package postgres

import (
	"context"

	"github.com/Traunin/review-assigner/internal/domain/entities"
	"github.com/Traunin/review-assigner/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *entities.User,
) error {
	_, err := r.db.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		UserID:    user.ID().String(),
		Username:  user.Username(),
		IsActive:  user.IsActive(),
	})

	return err
}

func (r *UserRepository) DeleteByID(
	ctx context.Context,
	id entities.UserID,
) error {
	return r.db.Queries.DeleteUser(ctx, id.String())
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id entities.UserID,
) (*entities.User, error) {
	user, err := r.db.Queries.GetUserByID(ctx, id.String())
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return entities.NewUser(
		entities.UserID(user.UserID),
		user.Username,
		user.IsActive,
	)
}
func (r *UserRepository) FindAll(
	ctx context.Context,
) ([]*entities.User, error) {
	users, err := r.db.Queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	userEntities := make([]*entities.User, len(users))
	for i, user := range users {
		userEntity, err := entities.NewUser(
			entities.UserID(user.UserID),
			user.Username,
			user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		userEntities[i] = userEntity
	}

	return userEntities, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	user *entities.User,
) error {
	return r.db.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:    user.ID().String(),
		Username:  user.Username(),
		IsActive:  user.IsActive(),
	})
}

func (r *UserRepository) GetActiveUsers(
	ctx context.Context,
) ([]*entities.User, error) {
	users, err := r.db.Queries.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	userEntities := make([]*entities.User, len(users))
	for i, user := range users {
		userEntity, err := entities.NewUser(
			entities.UserID(user.UserID),
			user.Username,
			user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		userEntities[i] = userEntity
	}

	return userEntities, nil
}
