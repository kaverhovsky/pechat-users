package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
	"pechat-users/internal/domain/model"
)

type Repository interface {
	GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error)
	GetAllUsersViews(ctx context.Context) ([]*model.UserView, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, ID uuid.UUID, opts *model.UserUpdateOpts) error
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (s *UseCase) GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, ID)
	if err != nil {
		logger.Error(ctx, "failed to get user by id",
			zap.NamedError("err", err),
			zap.String("userId", ID.String()))
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	logger.Debug(ctx, "got user by id",
		zap.String("userId", ID.String()))
	return user, nil
}

func (s *UseCase) GetAllUsersViews(ctx context.Context) ([]*model.UserView, error) {
	users, err := s.repo.GetAllUsersViews(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get all users views", zap.NamedError("err", err))
		return nil, fmt.Errorf("failed to get all users views: %w", err)
	}

	logger.Debug(ctx, "got all users views")

	return users, nil
}

func (s *UseCase) CreateUser(ctx context.Context, opts *model.UserCreateOpts) error {
	ID := uuid.New()

	// TODO hash password
	// passwordHash := hash(opts.Password)

	newUser := &model.User{
		ID:           ID,
		Nickname:     opts.Nickname,
		PasswordHash: opts.Password, // TODO put hash
		Firstname:    &opts.Firstname,
		Lastname:     &opts.Lastname,
		Email:        opts.Email,
		Bio:          &opts.Bio,
	}

	err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Error(ctx, "failed to create user", zap.NamedError("err", err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.Debug(ctx, "updated user by id",
		zap.String("userId", ID.String()))

	return nil
}

func (s *UseCase) UpdateUser(ctx context.Context, ID uuid.UUID, opts *model.UserUpdateOpts) error {
	if err := s.repo.UpdateUser(ctx, ID, opts); err != nil {
		logger.Error(ctx, "failed to update user by id",
			zap.NamedError("err", err),
			zap.String("userId", ID.String()))
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.Debug(ctx, "updated user by id",
		zap.String("userId", ID.String()))

	return nil
}

func (s *UseCase) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	if err := s.repo.DeleteUserByID(ctx, ID); err != nil {
		logger.Error(ctx, "failed to delete user by id",
			zap.NamedError("err", err),
			zap.String("userId", ID.String()))
		return fmt.Errorf("failed to delete user by id: %w", err)
	}

	logger.Debug(ctx, "deleted user by id",
		zap.String("userId", ID.String()))

	return nil
}
