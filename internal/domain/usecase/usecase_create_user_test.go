package usecase

import (
	"context"
	"fmt"
	"github.com/kaverhovsky/pechat-lib/logger"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pechat-users/internal/domain/model"
	"pechat-users/internal/domain/usecase/mocks"
	"testing"
)

// custom matcher for model.User
type createUserMatcher struct {
	expected *model.User
}

func (m *createUserMatcher) Matches(x interface{}) bool {
	got, ok := x.(*model.User)
	if !ok {
		return false
	}

	// check all fields
	// except ID - it is generated in CreateUser
	// except PasswordHash - hash is different after every generation
	return m.expected.Nickname == got.Nickname &&
		m.expected.Email == got.Email &&
		*m.expected.Firstname == *got.Firstname &&
		*m.expected.Lastname == *got.Lastname &&
		*m.expected.Bio == *got.Bio
}

func (m *createUserMatcher) String() string {
	str := fmt.Sprintf(`expected user struct (ignored ID):
	Nickname: %s,
	Email: %s,
	Firstname: %s,
	Lastname: %s,
	Bio: %s`, m.expected.Nickname, m.expected.Email, *m.expected.Firstname, *m.expected.Lastname, *m.expected.Bio)
	return str
}

func matchCreatedUser(expected *model.User) gomock.Matcher {
	return &createUserMatcher{expected}
}

func TestCreateUser(t *testing.T) {
	type TestCase struct {
		name           string
		userCreateOpts *model.UserCreateOpts
		mockSetup      func(repo *mocks.MockRepository)
		expectedErr    bool
	}

	testCases := []*TestCase{
		{
			name: "success",
			userCreateOpts: &model.UserCreateOpts{
				Nickname:  "jojofan",
				Password:  "abcd1234",
				Firstname: "Giorno",
				Lastname:  "Giovanni",
				Email:     "giogio@example.com",
				Bio:       "hello i'm a giogio",
			},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateUser(context.Background(),
					matchCreatedUser(&model.User{
						Nickname:     "jojofan",
						PasswordHash: "abcd1234",
						Firstname:    lo.ToPtr("Giorno"),
						Lastname:     lo.ToPtr("Giovanni"),
						Email:        "giogio@example.com",
						Bio:          lo.ToPtr("hello i'm a giogio"),
					})).Return(nil)
			},
			expectedErr: false,
		},
		//{
		//	name: "success",
		//	userCreateOpts: &model.UserCreateOpts{
		//		Nickname:  "jojofan",
		//		Password:  "abcd1234",
		//		Firstname: "Giorno",
		//		Lastname:  "Giovanni",
		//		Email:     "giogio@example.com",
		//		Bio:       "hello i'm a giogio",
		//	},
		//	mockSetup: func(repo *mocks.MockRepository) {
		//		repo.EXPECT().CreateUser(context.Background(),
		//			matchCreatedUser(&model.User{
		//				Nickname:     "jojofan",
		//				PasswordHash: "abcd1234",
		//				Firstname:    lo.ToPtr("Giorno"),
		//				Lastname:     lo.ToPtr("Giovanni"),
		//				Email:        "giogio@example.com",
		//				Bio:          lo.ToPtr("hello i'm a giogio"),
		//			})).Return(nil)
		//	},
		//	expectedErr: false,
		//},

		// What else to test?
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			controller := gomock.NewController(t)

			mockRepo := mocks.NewMockRepository(controller)

			tt.mockSetup(mockRepo)

			logger.SetupLogger("production", "info")
			uc := NewUseCase(mockRepo)

			err := uc.CreateUser(context.Background(), tt.userCreateOpts)
			if tt.expectedErr {
				require.Error(t, err, "must be an error")
			} else {
				require.NoError(t, err, "must be no error")
			}

		})
	}
}
