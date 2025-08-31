package get_all_users_views_handler

import (
	"context"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
	"net/http"
	"pechat-users/internal/domain/model"
	"pechat-users/pkg/http_helpers"
)

type useCase interface {
	GetAllUsersViews(ctx context.Context) ([]*model.UserView, error)
}

type Handler struct {
	uc useCase
}

func New(useCase useCase) *Handler {
	return &Handler{uc: useCase}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.uc.GetAllUsersViews(ctx)
	if err != nil {
		logger.Info(ctx, "failed to get all users views", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http_helpers.WriteSuccess(w, http.StatusOK, users)
}
