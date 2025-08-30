package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/kaverhovsky/pechat-lib/logger"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pechat-users/internal/domain/model"
	"pechat-users/internal/domain/usecase/mocks"
	"testing"
)

func TestGetAllUsersViews(t *testing.T) {
	type TestCase struct {
		name          string
		mockSetup     func(repo *mocks.MockRepository)
		expectedViews []*model.UserView
		expectedErr   bool
	}

	id1 := uuid.New()
	id2 := uuid.New()
	testCases := []*TestCase{
		{
			name: "success",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetAllUsersViews(context.Background()).
					Return([]*model.UserView{
						{
							ID:        id1,
							Nickname:  "narutofan",
							Firstname: lo.ToPtr("Sasuke"),
							Lastname:  lo.ToPtr("Uchiha"),
						},
						{
							ID:        id2,
							Nickname:  "jojofan",
							Firstname: lo.ToPtr("Giorno"),
							Lastname:  lo.ToPtr("Giovanni"),
						},
					},
						nil)
			},
			expectedViews: []*model.UserView{
				{
					ID:        id1,
					Nickname:  "narutofan",
					Firstname: lo.ToPtr("Sasuke"),
					Lastname:  lo.ToPtr("Uchiha"),
				},
				{
					ID:        id2,
					Nickname:  "jojofan",
					Firstname: lo.ToPtr("Giorno"),
					Lastname:  lo.ToPtr("Giovanni"),
				},
			},
			expectedErr: false,
		},
		{
			name: "no users",
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetAllUsersViews(context.Background()).
					Return([]*model.UserView{},
						nil)
			},
			expectedViews: []*model.UserView{},
			expectedErr:   false,
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

			got, err := uc.GetAllUsersViews(context.Background())
			if tt.expectedErr {
				require.Error(t, err, "must be an error")
			} else {
				require.NoError(t, err, "must be no error")
			}

			require.EqualValues(t, tt.expectedViews, got, "users must be equal")
		})
	}
}
