package controllers

import (
	"context"
	"fmt"
	controllers_test "github.com/alexsibrin/runbot-auth/internal/api/controllers/mocks"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	usecase := controllers_test.NewMockIAccountUsecase(ctrl)
	securer := controllers_test.NewMockISecurer(ctrl)

	testcases := []struct {
		name        string
		in          *models.SignUp
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Valid sign up",
			in: &models.SignUp{
				Email:    "test@test.ru",
				Password: "strongpswd",
				Name:     "SomeName",
			},
			setupMocks: func() {
				usecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&entities.Account{}, nil)
				securer.EXPECT().AccessToken(gomock.Any()).Return("atoken", nil)
				securer.EXPECT().RefreshToken(gomock.Any()).Return("rtoken", nil)
			},
			expectedErr: nil,
		},
		{
			name: "Short password",
			in: &models.SignUp{
				Email:    "test@test.ru",
				Password: "pswd",
				Name:     "SomeName",
			},
			setupMocks:  func() {},
			expectedErr: validators.ErrPasswordIsTooShort,
		},
		{
			name: "Short email",
			in: &models.SignUp{
				Email:    "test",
				Password: "strongpswd",
				Name:     "SomeName",
			},
			setupMocks:  func() {},
			expectedErr: validators.ErrEmailIsTooShort,
		},
		{
			name: "Wrong email",
			in: &models.SignUp{
				Email:    "test33",
				Password: "strongpswd",
				Name:     "SomeName",
			},
			setupMocks:  func() {},
			expectedErr: validators.ErrEmailFormatIsNotCorrect,
		},
		{
			name: "Short name",
			in: &models.SignUp{
				Email:    "test@test.ru",
				Password: "strongpswd",
				Name:     "sv",
			},
			setupMocks:  func() {},
			expectedErr: validators.ErrNameIsTooShort,
		},
		{
			name: "Wrong name",
			in: &models.SignUp{
				Email:    "test@test.ru",
				Password: "strongpswd",
				Name:     "SomeName!!!",
			},
			setupMocks:  func() {},
			expectedErr: validators.ErrNameFormatIsNotCorrect,
		},
		{
			name: "Account is already exist",
			in: &models.SignUp{
				Email:    "test@test.ru",
				Password: "strongpswd",
				Name:     "SomeName",
			},
			setupMocks: func() {
				usecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrAccountAlreadyExist)
			},
			expectedErr: usecases.ErrAccountAlreadyExist,
		},
		{
			name: "Other account usecase error",
			in: &models.SignUp{
				Email:    "test@test.ru",
				Password: "strongpswd",
				Name:     "SomeName",
			},
			setupMocks: func() {
				usecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error"))
			},
			expectedErr: fmt.Errorf("some error"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{
				usecase: usecase,
				securer: securer,
			}

			acc, err := account.SignUp(ctx, tc.in)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.IsType(t, &models.SignUpResponse{}, acc)
				assert.NoError(t, err)
			}

		})
	}
}
