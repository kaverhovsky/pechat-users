package create_user_handler

import (
	"context"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kaverhovsky/pechat-lib/logger"
	"go.uber.org/zap"
	"net/http"
	"pechat-users/internal/domain/model"
	"pechat-users/pkg/http_helpers"
)

type useCase interface {
	CreateUser(ctx context.Context, opts *model.UserCreateOpts) error
}

type Handler struct {
	uc useCase
}

func New(uc useCase) *Handler {
	return &Handler{uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createUserRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode create user request", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Nickname, validation.Required, validation.Length(3, 100)),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required, validation.Length(6, 100)),
		validation.Field(&req.Firstname, validation.Length(1, 100)),
		validation.Field(&req.Lastname, validation.Length(1, 100)),
		validation.Field(&req.Bio, validation.Length(1, 300)),
	); err != nil {
		logger.Error(ctx, "failed to validate create user request", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.uc.CreateUser(ctx, &model.UserCreateOpts{
		Nickname:  req.Nickname,
		Email:     req.Email,
		Password:  req.Password,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Bio:       req.Bio,
	}); err != nil {
		logger.Error(ctx, "failed to create user", zap.NamedError("err", err))
		http_helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http_helpers.WriteSuccess(w, http.StatusCreated, nil)
}
