package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"math/big"
	"pechat-users/internal/domain/model"
	"pechat-users/internal/pkg/config"
	"testing"
	"time"
)

//func TestMain(m *testing.M) {
//
//}

func setup(ctx context.Context) (*PostgresRepository, *big.Int, error) {
	conf, err := config.LoadConfig("../../../configs/local.yaml")
	if err != nil {
		return nil, &big.Int{}, fmt.Errorf("failed to load config: %w", err)
	}

	poolConf, err := pgxpool.ParseConfig(conf.Postgres.DSN)
	if err != nil {
		return nil, &big.Int{}, fmt.Errorf("failed to parse pgxpool config from dsn: %w", err)
	}
	// для сохранения состояния между запросами
	poolConf.MinConns = 1
	poolConf.MaxConns = 1
	schemaID, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, nil, err
	}
	poolConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err = conn.Exec(ctx, fmt.Sprintf("create schema if not exists test%v;", schemaID))
		if err != nil {
			return err
		}

		_, err = conn.Exec(ctx, fmt.Sprintf("set search_path to test%s;", schemaID.String()))
		if err != nil {
			return err
		}
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, &big.Int{}, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	//schemaID, err := createAndSetSchema(ctx, pool)
	//if err != nil {
	//	return nil, &big.Int{}, fmt.Errorf("failed to create and set schema: %w", err)
	//}

	if err := applyMigrations("../../../migrations", pool); err != nil {
		return nil, &big.Int{}, fmt.Errorf("failed to apply migrations: %w", err)
	}

	repo := NewPostgresRepositoryWithPool(pool)

	return repo, schemaID, nil
}

func cleanup(ctx context.Context, pool *pgxpool.Pool, schemaID *big.Int) (err error) {
	_, err = pool.Exec(ctx, fmt.Sprintf("drop schema test%s cascade", schemaID.String()))
	pool.Close()
	return
}

func TestGetUserByID_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaID, err := setup(ctx)
	//repo, _, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	defer func() {
		err := cleanup(ctx, repo.pool, schemaID)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	//repo, err := NewPostgresRepository(ctx, "postgres://postgres:password@localhost:5432/pechat_users?sslmode=disable")
	//require.NoError(t, err, "setting up postgres repository")
	// test

	q := `insert into users 
(id, nickname, password_hash, firstname, lastname, email, bio)
values
($1, $2, $3, $4, $5, $6, $7)`
	u := &model.User{
		ID:           uuid.New(),
		Nickname:     "jojofan",
		PasswordHash: "abcd1234",
		Firstname:    lo.ToPtr("Giorno"),
		Lastname:     lo.ToPtr("Giovanni"),
		Email:        "giogio@example.com",
		Bio:          lo.ToPtr("hello i'm a giogio"),
	}

	_, err = repo.pool.Exec(ctx, q, u.ID, u.Nickname, u.PasswordHash, u.Firstname, u.Lastname, u.Email, u.Bio)
	require.NoError(t, err, "should be no error while inserting test user to repo")

	got, err := repo.GetUserByID(ctx, u.ID)
	require.NoError(t, err, "got error while getting user by id")

	require.Equal(t, got, u, "expected user and acquired user are not the same")
}
