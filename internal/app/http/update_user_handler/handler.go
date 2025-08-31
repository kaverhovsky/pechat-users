package update_user_handler

import (
	"context"
	"encoding/json"
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
	UpdateUser(ctx context.Context, ID uuid.UUID, opts *model.UserUpdateOpts) error
}

type Handler struct {
	uc useCase
}

func New(uc useCase) *Handler {
	return &Handler{uc}
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

	var req UpdateUserRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode update user request", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Nickname, validation.NilOrNotEmpty, validation.Length(3, 100)),
		validation.Field(&req.Firstname, validation.NilOrNotEmpty, validation.Length(1, 100)),
		validation.Field(&req.Lastname, validation.NilOrNotEmpty, validation.Length(1, 100)),
		validation.Field(&req.Bio, validation.NilOrNotEmpty, validation.Length(1, 300)),
	); err != nil {
		logger.Error(ctx, "failed to validate update user request", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.uc.UpdateUser(ctx, userID, &model.UserUpdateOpts{
		Nickname:  req.Nickname,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Bio:       req.Bio,
	}); err != nil {
		logger.Error(ctx, "failed to update user", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http_helpers.WriteSuccess(w, http.StatusCreated, nil)
}
