package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kaverhovsky/chat/internal/domain/model"
)

type Repository interface {
	GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error)
	GetAllUsersViews(ctx context.Context) ([]*model.UserView, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, opts *model.UserUpdateOpts) error
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
}

type UseCase struct {
	repo Repository
}

func NewService(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (s *UseCase) CreateUser(ctx context.Context, opts *model.UserCreateOpts) error {
	ID := uuid.New()

	// TODO hash password
	// passwordHash := hash(opts.Password)

	newUser := &model.User{
		ID:           ID,
		Nickname:     opts.Nickname,
		PasswordHash: opts.Password, // TODO put hash
		Firstname:    opts.Firstname,
		Lastname:     opts.Lastname,
		Email:        opts.Email,
		Bio:          opts.Bio,
	}

	err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil

}

