package dto

import "github.com/Traunin/review-assigner/internal/domain/entities"

type CreateUserCmd struct {
	UserID   entities.UserID
	Username string
	TeamID   *int64
	IsActive *bool
}

type UpdateUserCmd struct {
	Username string
	TeamID   *int64
	IsActive bool
}

type UserDTO struct {
	UserID   entities.UserID
	Username string
	IsActive bool
	TeamID   *int64
}
