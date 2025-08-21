package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"math/big"
)

func createAndSetSchema(ctx context.Context, pool *pgxpool.Pool) (*big.Int, error) {
	schemaID, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, fmt.Sprintf("create schema test%v;", schemaID))
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, fmt.Sprintf("set search_path to test%s;", schemaID.String()))
	if err != nil {
		return nil, err
	}
	return schemaID, nil
}

func applyMigrations(migrationPath string, pool *pgxpool.Pool) (err error) {
	if err = goose.SetDialect("postgres"); err != nil {
		return
	}

	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		err = db.Close()
	}()

	if err = goose.Up(db, migrationPath); err != nil {
		return
	}

	return
}
