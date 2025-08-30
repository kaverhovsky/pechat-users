package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kaverhovsky/pechat-lib/logger"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pechat-users/internal/domain/model"
	"pechat-users/internal/domain/usecase/mocks"
	"testing"
)

// custom matcher for model.User
type updateUserOptsMatcher struct {
	expected *model.UserUpdateOpts
}

func (m *updateUserOptsMatcher) Matches(x interface{}) bool {
	got, ok := x.(*model.UserUpdateOpts)
	if !ok {
		return false
	}

	// check all fields except ID, because it is generated in CreateUser
	return *m.expected.Nickname == *got.Nickname &&
		*m.expected.Firstname == *got.Firstname &&
		*m.expected.Lastname == *got.Lastname &&
		*m.expected.Bio == *got.Bio
}

func (m *updateUserOptsMatcher) String() string {
	str := fmt.Sprintf(`expected user struct (ignored ID and PasswordHash):
	Nickname: %s,
	Firstname: %s,
	Lastname: %s,
	Bio: %s`, *m.expected.Nickname, *m.expected.Firstname, *m.expected.Lastname, *m.expected.Bio)
	return str
}

func matchUpdatedUser(expected *model.UserUpdateOpts) gomock.Matcher {
	return &updateUserOptsMatcher{expected}
}

func TestUpdateUser(t *testing.T) {
	type TestCase struct {
		name           string
		ID             uuid.UUID
		userUpdateOpts *model.UserUpdateOpts
		mockSetup      func(repo *mocks.MockRepository)
		expectedUser   *model.User
		expectedErr    bool
	}

	id := uuid.New()
	testCases := []*TestCase{
		{
			name: "success",
			ID:   id,
			userUpdateOpts: &model.UserUpdateOpts{
				Nickname:  lo.ToPtr("jojofan"),
				Firstname: lo.ToPtr("Giorno"),
				Lastname:  lo.ToPtr("Giovanni"),
				Bio:       lo.ToPtr("hello i'm a giogio"),
			},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().UpdateUser(context.Background(),
					id,
					matchUpdatedUser(&model.UserUpdateOpts{
						Nickname:  lo.ToPtr("jojofan"),
						Firstname: lo.ToPtr("Giorno"),
						Lastname:  lo.ToPtr("Giovanni"),
						Bio:       lo.ToPtr("hello i'm a giogio"),
					})).Return(nil)
			},
			expectedErr: false,
		},

		// What else to test?
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)

			mockRepo := mocks.NewMockRepository(controller)

			tt.mockSetup(mockRepo)

			logger.SetupLogger("production", "info")
			uc := NewUseCase(mockRepo)

			err := uc.UpdateUser(context.Background(), tt.ID, tt.userUpdateOpts)
			if tt.expectedErr {
				require.Error(t, err, "must be an error")
			} else {
				require.NoError(t, err, "must be no error")
			}

		})
	}
}
