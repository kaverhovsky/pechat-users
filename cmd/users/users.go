package main

import (
	"context"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pgrepo, pgErr := repository.NewPostgresRepository(ctx, c.Postgres.DSN)
	if pgErr != nil {
		logger.Logger().Error("failed to create postgres repository", zap.NamedError("err", pgErr))
	}

	_ = usecase.NewUseCase(pgrepo)
}
