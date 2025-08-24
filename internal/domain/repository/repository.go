package repository

import (
	"context"
	"fmt"
	j "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"pechat-users/internal/domain/model"
	j_model "pechat-users/internal/pkg/jet/pechat_users/public/model"
	t "pechat-users/internal/pkg/jet/pechat_users/public/table"
)

type PostgresRepository struct {
	pool  *pgxpool.Pool
	table *t.UsersTable
}

type RepoOption func(repo *PostgresRepository) error

func WithSchema(schemaName string) RepoOption {
	return func(repo *PostgresRepository) error {
		repo.table = t.Users.FromSchema(schemaName)
		return nil
	}
}

func NewPostgresRepository(ctx context.Context, dsn string, opts ...RepoOption) (*PostgresRepository, error) {
	// TODO указать другие параметры (размер пула и пр.)
	c, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgxpool config from dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	repo := &PostgresRepository{pool: pool}
	repo.table = t.Users

	for _, opt := range opts {
		if err := opt(repo); err != nil {
			return nil, fmt.Errorf("failed to apply repo option: %w", err)
		}
	}

	return repo, nil
}

func NewPostgresRepositoryWithPool(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (p *PostgresRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	stmt := j.SELECT(p.table.AllColumns).FROM(p.table).WHERE(p.table.ID.EQ(j.UUID(ID)))

	q, args := stmt.Sql()

	row := p.pool.QueryRow(ctx, q, args...)

	var u j_model.Users
	if err := row.Scan(&u.ID, &u.Nickname, &u.PasswordHash, &u.Firstname, &u.Lastname, &u.Email, &u.Bio); err != nil {
		return nil, fmt.Errorf("failed to scan user field values: %w", err)
	}

	return &model.User{
		ID:           u.ID,
		Nickname:     u.Nickname,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		Firstname:    u.Firstname,
		Lastname:     u.Lastname,
		Bio:          u.Bio,
	}, nil
}

func (p *PostgresRepository) GetAllUsersViews(ctx context.Context) ([]*model.UserView, error) {
	// TODO limit
	stmt := j.SELECT(p.table.ID, p.table.Nickname, p.table.Firstname, p.table.Lastname).FROM(p.table)
	q, args := stmt.Sql()

	var us []j_model.Users

	rows, err := p.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users views: %w", err)
	}

	for rows.Next() {
		var u j_model.Users
		if err := rows.Scan(&u.ID, &u.Nickname, &u.Firstname, &u.Lastname); err != nil {
			return nil, fmt.Errorf("failed to scan row to jet user model: %w", err)
		}
		us = append(us, u)
	}

	views := make([]*model.UserView, 0)
	for _, u := range us {
		views = append(views, &model.UserView{
			ID:        u.ID,
			NickName:  u.Nickname,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
		})
	}

	return views, nil
}

func (p *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	u := j_model.Users{
		ID:           user.ID,
		Nickname:     user.Nickname,
		PasswordHash: user.PasswordHash,
		Firstname:    user.Firstname,
		Lastname:     user.Lastname,
		Email:        user.Email,
		Bio:          user.Bio,
	}

	stmt := p.table.INSERT(
		p.table.ID,
		p.table.Nickname,
		p.table.PasswordHash,
		p.table.Firstname,
		p.table.Lastname,
		p.table.Email,
		p.table.Bio,
	).MODEL(u)

	q, args := stmt.Sql()
	if _, err := p.pool.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to execute insert query for user: %w", err)
	}

	return nil
}

// UpdateUser TODO returning?
func (p *PostgresRepository) UpdateUser(ctx context.Context, ID uuid.UUID, opts *model.UserUpdateOpts) error {
	stmt := p.table.UPDATE()
	if opts.Nickname != nil {
		stmt = stmt.SET(p.table.Nickname.SET(j.String(*opts.Nickname)))
	}
	if opts.Firstname != nil {
		stmt = stmt.SET(p.table.Firstname.SET(j.String(*opts.Firstname)))
	}
	if opts.Lastname != nil {
		stmt = stmt.SET(p.table.Lastname.SET(j.String(*opts.Lastname)))
	}
	if opts.Bio != nil {
		stmt = stmt.SET(p.table.Bio.SET(j.String(*opts.Bio)))
	}
	stmt = stmt.WHERE(p.table.ID.EQ(j.UUID(ID)))

	q, args := stmt.Sql()

	if _, err := p.pool.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

func (p *PostgresRepository) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	stmt := p.table.DELETE().WHERE(p.table.ID.EQ(j.UUID(ID)))
	q, args := stmt.Sql()

	if _, err := p.pool.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("failed to execute delete by id query for user: %w", err)
	}

	return nil
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}
