package get_user_by_id_handler

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
	"net/http"
	"pechat-users/internal/domain/model"
	"pechat-users/pkg/http_helpers"
)

type useCase interface {
	GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error)
}

type Handler struct {
	uc useCase
}

func New(useCase useCase) *Handler {
	return &Handler{uc: useCase}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := r.PathValue("userID")
	if err := validation.Validate(userIDStr, validation.Required, is.UUID); err != nil {
		logger.Error(ctx, "validation failed for userID", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := uuid.MustParse(userIDStr)
	user, err := h.uc.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "failed to get user by id", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http_helpers.WriteSuccess(w, http.StatusOK, user)
}
