package main

import (
	"context"
	"github.com/kaverhovsky/pechat-lib/logger"
	"log"
	"pechat-users/internal/domain/repository"
	"pechat-users/internal/domain/usecase"
	"pechat-users/internal/pkg/config"
	"time"
)

func main() {
	// TODO use flag
	c, err := config.LoadConfig("./configs/local")
	if err != nil {
		panic(err)
	}

	logger.SetupLogger(c.Logger.Mode, c.Logger.Level)

	pgrepo := repository.NewPostgresRepository()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := pgrepo.Init(ctx, c.Postgres.DSN); err != nil {
		// TODO log error and exit gracefully
		log.Fatal("err")
	}

	_ = usecase.NewUseCase(pgrepo)
}
