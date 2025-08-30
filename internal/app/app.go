package app

import (
	"context"
	"errors"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
	"net/http"
	"pechat-users/internal/app/http/create_user_handler"
	"pechat-users/internal/app/http/status_handler"
	"pechat-users/internal/domain/usecase"
	"pechat-users/internal/pkg/config"
)

type App struct {
	conf   *config.Config
	server *http.Server
}

func NewApp(conf *config.Config, uc *usecase.UseCase) *App {
	handler := applyHandlers(uc)

	app := &App{
		conf: conf,
		// TODO setup server configuration
		server: &http.Server{
			Addr:    conf.HTTP.Listen,
			Handler: handler,
		},
	}

	return app
}

func (app *App) Run() {
	logger.Logger().Info("starting http server...")
	if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Logger().Error("unexpected error from ListenAndServe http server",
			zap.NamedError("err", err))
	}
}

func applyHandlers(uc *usecase.UseCase) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /status", status_handler.New())
	mux.Handle("POST /api/v1/users", create_user_handler.New(uc))

	return mux
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
