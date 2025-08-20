package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaverhovsky/chat/internal/domain/model"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (p *PostgresRepository) Init(ctx context.Context, dsn string) error {
	c, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse pgxpool config from dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}
	p.pool = pool

	return nil
}

func (p *PostgresRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) GetAllUsersViews(ctx context.Context) ([]*model.UserView, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) UpdateUser(ctx context.Context, opts *model.UserUpdateOpts) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}
