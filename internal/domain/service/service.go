package service

import (
	"github.com/google/uuid"
	"github.com/kaverhovsky/chat/internal/domain/model"
)

type Repository interface {
	GetUserByID(ID uuid.UUID) (*model.User, error)
	GetAllUsersViews() ([]*model.UserView, error)
	CreateUser(user *model.User) error
	UpdateUser(opts *model.UserUpdateOpts) error
	DeleteUserByID(ID uuid.UUID) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
