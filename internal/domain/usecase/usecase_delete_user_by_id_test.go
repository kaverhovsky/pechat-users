package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/kaverhovsky/pechat-lib/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pechat-users/internal/domain/usecase/mocks"
	"testing"
)

func TestDeleteUserByID(t *testing.T) {
	type TestCase struct {
		name        string
		ID          uuid.UUID
		mockSetup   func(repo *mocks.MockRepository)
		expectedErr bool
	}

	id := uuid.New()

	testCases := []*TestCase{
		{
			name: "success",
			ID:   id,
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().DeleteUserByID(context.Background(), id).
					Return(nil)
			},
			expectedErr: false,
		},
		{
			name: "nothing to delete",
			ID:   id,
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().DeleteUserByID(context.Background(), id).
					Return(nil)
			},
			expectedErr: false,
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

			err := uc.DeleteUserByID(context.Background(), tt.ID)
			if tt.expectedErr {
				require.Error(t, err, "must be an error")
			} else {
				require.NoError(t, err, "must be no error")
			}
		})
	}
}
