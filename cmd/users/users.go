package main

import (
	"context"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"pechat-users/internal/app"
	"pechat-users/internal/domain/repository"
	"pechat-users/internal/domain/usecase"
	"pechat-users/internal/pkg/config"
	"syscall"
	"time"
)

func main() {
	// TODO use flag
	c, err := config.LoadConfig("./configs/local.yaml")
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

	uc := usecase.NewUseCase(pgrepo)

	httpApp := app.NewApp(c, uc)

	go httpApp.Run()

	defer func() {
		if err := httpApp.Shutdown(ctx); err != nil {
			logger.Logger().Error("failed to shutdown http app", zap.Error(err))
		}
	}()

	logger.Logger().Info("started users service")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	logger.Logger().Info("shutting down service...")
}
