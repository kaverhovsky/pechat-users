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
	pool *pgxpool.Pool
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (p *PostgresRepository) Init(ctx context.Context, dsn string) error {
	// TODO указать другие параметры (размер пула и пр.)
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
	stmt := j.SELECT(t.Users.AllColumns).FROM(t.Users).WHERE(t.Users.ID.EQ(j.Text(ID.String())))
	q, args := stmt.Sql()

	row := p.pool.QueryRow(ctx, q, args)

	var u j_model.Users
	if err := row.Scan(&u.ID, &u.Nickname, &u.PasswordHash, u.Firstname, u.Lastname, &u.Email, u.Bio); err != nil {
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
	stmt := j.SELECT(t.Users.ID, t.Users.Nickname, t.Users.Firstname, t.Users.Lastname).FROM(t.Users)
	q, args := stmt.Sql()

	var us []j_model.Users

	rows, err := p.pool.Query(ctx, q, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query users views: %w", err)
	}
	
	for rows.Next() {
		var u j_model.Users
		if err := rows.Scan(&u.ID, &u.Nickname, u.Firstname, u.Lastname); err != nil {
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
