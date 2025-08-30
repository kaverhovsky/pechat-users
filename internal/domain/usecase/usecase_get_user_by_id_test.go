package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/kaverhovsky/pechat-lib/logger"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pechat-users/internal/domain/model"
	"pechat-users/internal/domain/repository"
	"pechat-users/internal/domain/usecase/mocks"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	type TestCase struct {
		name         string
		ID           uuid.UUID
		mockSetup    func(repo *mocks.MockRepository)
		expectedUser *model.User
		expectedErr  bool
	}

	id := uuid.New()

	testCases := []*TestCase{
		{
			name: "success",
			ID:   id,
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetUserByID(context.Background(), id).
					Return(&model.User{
						ID:           id,
						Nickname:     "jojofan",
						PasswordHash: "abcd1234",
						Firstname:    lo.ToPtr("Giorno"),
						Lastname:     lo.ToPtr("Giovanni"),
						Email:        "giogio@example.com",
						Bio:          lo.ToPtr("hello i'm a giogio"),
					},
						nil)
			},
			expectedUser: &model.User{
				ID:           id,
				Nickname:     "jojofan",
				PasswordHash: "abcd1234",
				Firstname:    lo.ToPtr("Giorno"),
				Lastname:     lo.ToPtr("Giovanni"),
				Email:        "giogio@example.com",
				Bio:          lo.ToPtr("hello i'm a giogio"),
			},
			expectedErr: false,
		},
		{
			name: "not found",
			ID:   id,
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetUserByID(context.Background(), id).
					Return(nil, repository.ErrNotFound)
			},
			expectedUser: nil,
			expectedErr:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)

			mockRepo := mocks.NewMockRepository(controller)

			tt.mockSetup(mockRepo)

			logger.SetupLogger("production", "info")
			uc := NewUseCase(mockRepo)

			got, err := uc.GetUserByID(context.Background(), tt.ID)
			if tt.expectedErr {
				require.Error(t, err, "must be an error")
			} else {
				require.NoError(t, err, "must be no error")
			}

			require.Equal(t, tt.expectedUser, got, "users must be equal")
		})
	}
}
