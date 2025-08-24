package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
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

// setup returns repo, schemaName for consequent cleanup and an error
func setup(ctx context.Context) (*PostgresRepository, string, error) {
	conf, err := config.LoadConfig("../../../configs/local.yaml")
	if err != nil {
		return nil, "", fmt.Errorf("failed to load config: %w", err)
	}

	poolConf, err := pgxpool.ParseConfig(conf.Postgres.DSN)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse pgxpool config from dsn: %w", err)
	}

	schemaID, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, "", err
	}
	schemaName := "test" + schemaID.String()

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// test schema creation
	_, err = pool.Exec(ctx, fmt.Sprintf("create schema %s", schemaName))
	if err != nil {
		return nil, "", err
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("set search_path to %s", schemaName))
	if err != nil {
		return nil, "", err
	}

	if err := applyMigrations("../../../migrations", db); err != nil {
		return nil, "", fmt.Errorf("failed to apply migrations: %w", err)
	}

	//repo := NewPostgresRepositoryWithPool(pool)
	repo, err := NewPostgresRepository(ctx, conf.Postgres.DSN, WithSchema(schemaName))
	if err != nil {
		return nil, "", err
	}

	return repo, schemaName, nil
}

func cleanup(ctx context.Context, pool *pgxpool.Pool, schemaName string) (err error) {
	_, err = pool.Exec(ctx, fmt.Sprintf("drop schema %s cascade", schemaName))
	pool.Close()
	return
}

func TestGetUserByID_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code
	u := &model.User{
		ID:           uuid.New(),
		Nickname:     "jojofan",
		PasswordHash: "abcd1234",
		Firstname:    lo.ToPtr("Giorno"),
		Lastname:     lo.ToPtr("Giovanni"),
		Email:        "giogio@example.com",
		Bio:          lo.ToPtr("hello i'm a giogio"),
	}

	q, args := repo.table.INSERT(repo.table.AllColumns).VALUES(u.ID, u.Nickname, u.PasswordHash, u.Firstname, u.Lastname, u.Email, u.Bio).Sql()

	_, err = repo.pool.Exec(ctx, q, args...)
	require.NoError(t, err, "should be no error while inserting test user to repo")

	got, err := repo.GetUserByID(ctx, u.ID)
	require.NoError(t, err, "got error while getting user by id")

	require.Equal(t, got, u, "expected user and acquired user are not the same")
}

func TestGetUserByID_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code
	_, err = repo.GetUserByID(ctx, uuid.New())

	require.ErrorIs(t, err, pgx.ErrNoRows, "there must be not found error")
}

func TestGetAllUsersViews_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code
	u1 := &model.User{
		ID:           uuid.New(),
		Nickname:     "jojofan",
		PasswordHash: "abcd1234",
		Firstname:    lo.ToPtr("Giorno"),
		Lastname:     lo.ToPtr("Giovanni"),
		Email:        "giogio@example.com",
		Bio:          lo.ToPtr("hello i'm a giogio"),
	}
	u2 := &model.User{
		ID:           uuid.New(),
		Nickname:     "narutofan",
		PasswordHash: "qwerty0987",
		Firstname:    lo.ToPtr("Sasuke"),
		Lastname:     lo.ToPtr("Uchiha"),
		Email:        "amaterasu@example.com",
		Bio:          lo.ToPtr("hello i'm a sasuke"),
	}

	expected := []*model.UserView{
		{
			ID:        u1.ID,
			NickName:  u1.Nickname,
			Firstname: u1.Firstname,
			Lastname:  u1.Lastname,
		},
		{
			ID:        u2.ID,
			NickName:  u2.Nickname,
			Firstname: u2.Firstname,
			Lastname:  u2.Lastname,
		},
	}

	q, args := repo.table.INSERT(repo.table.AllColumns).
		VALUES(u1.ID, u1.Nickname, u1.PasswordHash, u1.Firstname, u1.Lastname, u1.Email, u1.Bio).
		VALUES(u2.ID, u2.Nickname, u2.PasswordHash, u2.Firstname, u2.Lastname, u2.Email, u2.Bio).Sql()

	_, err = repo.pool.Exec(ctx, q, args...)
	require.NoError(t, err, "should be no error while inserting test user to repo")

	got, err := repo.GetAllUsersViews(ctx)
	require.NoError(t, err, "got error while getting all users views")

	require.Equal(t, len(expected), len(got), "expected number of users views")
	require.EqualValues(t, got, expected, "expected user and acquired users views are not the same")
}

func TestGetAllUsersViews_NoRows(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code

	expected := make([]*model.UserView, 0)
	got, err := repo.GetAllUsersViews(ctx)
	require.NoError(t, err, "got error while getting all users views")

	require.Equal(t, len(expected), len(got), "expected number of users views")
	require.EqualValues(t, got, expected, "expected user and acquired users views are not the same")
}

func TestCreateUser_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code
	id := uuid.New()
	u := &model.User{
		ID:           id,
		Nickname:     "jojofan",
		PasswordHash: "abcd1234",
		Firstname:    lo.ToPtr("Giorno"),
		Lastname:     lo.ToPtr("Giovanni"),
		Email:        "giogio@example.com",
		Bio:          lo.ToPtr("hello i'm a giogio"),
	}
	err = repo.CreateUser(ctx, u)
	require.NoError(t, err, "should be no error while creating user")

	expected, err := repo.GetUserByID(ctx, id)
	require.NoError(t, err, "got error while getting user by id")

	require.Equal(t, u, expected, "expected user and acquired user are not the same")
}

// TODO should be an error?
func TestCreateUser_MissingFields(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code
	id := uuid.New()
	u := &model.User{
		ID:           id,
		Nickname:     "",
		PasswordHash: "",
		Firstname:    lo.ToPtr(""),
		Lastname:     lo.ToPtr(""),
		Email:        "",
		Bio:          lo.ToPtr(""),
	}
	err = repo.CreateUser(ctx, u)
	require.NoError(t, err, "should be no error while creating user")

	expected, err := repo.GetUserByID(ctx, id)
	require.NoError(t, err, "got error while getting user by id")

	require.Equal(t, u, expected, "expected user and acquired user are not the same")
}

func TestUpdateUser_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// setup repo
	repo, schemaName, err := setup(ctx)
	require.NoError(t, err, "setting up postgres repository")
	// cleanup on exit
	defer func() {
		err := cleanup(ctx, repo.pool, schemaName)
		require.NoError(t, err, "cleaning up postgres repository")
	}()

	// test code
	id := uuid.New()
	u := &model.User{
		ID:           id,
		Nickname:     "jojofan",
		PasswordHash: "abcd1234",
		Firstname:    lo.ToPtr("Giorno"),
		Lastname:     lo.ToPtr("Giovanni"),
		Email:        "giogio@example.com",
		Bio:          lo.ToPtr("hello i'm a giogio"),
	}
	q, args := repo.table.INSERT(repo.table.AllColumns).VALUES(u.ID, u.Nickname, u.PasswordHash, u.Firstname, u.Lastname, u.Email, u.Bio).Sql()

	_, err = repo.pool.Exec(ctx, q, args...)
	require.NoError(t, err, "should be no error while inserting test user to repo")

	newNickname := lo.ToPtr("narutofan")
	newFirstname := lo.ToPtr("Sasuke")
	newLastname := lo.ToPtr("Uchiha")
	newBio := lo.ToPtr("hello i'm a naruto fan")

	opts := &model.UserUpdateOpts{
		Nickname:  newNickname,
		Firstname: newFirstname,
		Lastname:  newLastname,
		Bio:       newBio,
	}

	u.Nickname = *newNickname
	u.Firstname = newFirstname
	u.Lastname = newLastname
	u.Bio = newBio

	// TODO returning
	err = repo.UpdateUser(ctx, id, opts)
	require.NoError(t, err, "should be no error while updating user")

	got, err := repo.GetUserByID(ctx, id)
	require.NoError(t, err, "got error while getting user by id")

	require.Equal(t, u, got, "expected user and updated user are not the same")
}
