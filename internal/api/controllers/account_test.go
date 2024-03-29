package controllers

import (
	"context"
	"database/sql"
	"fmt"
	controllers_test "github.com/alexsibrin/runbot-auth/internal/api/controllers/mocks"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
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
				assert.Errorf(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &models.SignUpResponse{}, acc)
			}

		})
	}
}

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := controllers_test.NewMockIAccountUsecase(ctrl)
	securer := controllers_test.NewMockISecurer(ctrl)

	ctx := context.TODO()

	testcases := []struct {
		name        string
		in          *models.SignIn
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Valid case",
			in: &models.SignIn{
				Email:    "test@test.ru",
				Password: "strongpswd",
			},
			setupMocks: func() {
				usecase.EXPECT().SignIn(ctx, gomock.Any(), gomock.Any()).Return(&entities.Account{}, nil)
				securer.EXPECT().AccessToken(gomock.Any()).Return("atoken", nil)
				securer.EXPECT().RefreshToken(gomock.Any()).Return("rtoken", nil)
			},
			expectedErr: nil,
		},
		{
			name: "Wrong password",
			in: &models.SignIn{
				Email:    "test@test.ru",
				Password: "somewrongpassword",
			},
			setupMocks: func() {
				usecase.EXPECT().SignIn(ctx, gomock.Any(), gomock.Any()).Return(nil, bcrypt.ErrMismatchedHashAndPassword)
			},
			expectedErr: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name: "Wrong email",
			in: &models.SignIn{
				Email:    "some@wrongemail.ru",
				Password: "strongpswd",
			},
			setupMocks: func() {
				usecase.EXPECT().SignIn(ctx, gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{
				usecase: usecase,
				securer: securer,
			}

			result, err := account.SignIn(ctx, tc.in)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &models.SignInResponse{}, result)
			}
		})
	}
}

func TestGetOneByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	usecase := controllers_test.NewMockIAccountUsecase(ctrl)

	testcases := []struct {
		name        string
		in          string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Valid case",
			in:   "some@validemail.ru",
			setupMocks: func() {
				usecase.EXPECT().GetOneByEmail(ctx, gomock.Any()).Return(&entities.Account{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Wrong email",
			in:   "some@wrongemail.ru",
			setupMocks: func() {
				usecase.EXPECT().GetOneByEmail(ctx, gomock.Any()).Return(nil, sql.ErrNoRows)
			},
			expectedErr: sql.ErrNoRows,
		},
		{
			name:        "Wrong email format",
			in:          "some@ru",
			setupMocks:  func() {},
			expectedErr: validators.ErrEmailFormatIsNotCorrect,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			account := &Account{
				usecase: usecase,
			}

			result, err := account.GetOneByEmail(ctx, tc.in)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &models.AccountGetModel{}, result)
			}

		})
	}
}

func TestGetOneByUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	usecase := controllers_test.NewMockIAccountUsecase(ctrl)

	testcases := []struct {
		name        string
		in          string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Valid case",
			in:   "validuuid",
			setupMocks: func() {
				usecase.EXPECT().GetOneByUUID(ctx, gomock.Any()).Return(&entities.Account{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Wrong uuid",
			in:   "wrongeuuid",
			setupMocks: func() {
				usecase.EXPECT().GetOneByUUID(ctx, gomock.Any()).Return(nil, sql.ErrNoRows)
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			account := &Account{
				usecase: usecase,
			}

			result, err := account.GetOneByUUID(ctx, tc.in)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &models.AccountGetModel{}, result)
			}

		})
	}
}

func TestRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securer := controllers_test.NewMockISecurer(ctrl)

	ctx := context.TODO()

	testCases := []struct {
		name        string
		in          string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Valid case",
			in:   "sometoken",
			setupMocks: func() {
				securer.EXPECT().Decrypt(gomock.Any()).Return(&entities.Account{}, nil)
				securer.EXPECT().RefreshToken(gomock.Any()).Return("newtoken", nil)
			},
			expectedErr: nil,
		},
		{
			name: "Wrong token signature",
			in:   "somewrongtoken",
			setupMocks: func() {
				securer.EXPECT().Decrypt(gomock.Any()).Return(nil, jwt.ErrTokenSignatureInvalid)
			},
			expectedErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			name: "Wrong token hash",
			in:   "somewrongtoken",
			setupMocks: func() {
				securer.EXPECT().Decrypt(gomock.Any()).Return(nil, jwt.ErrHashUnavailable)
			},
			expectedErr: jwt.ErrHashUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{
				securer: securer,
			}
			token, err := account.RefreshToken(ctx, tc.in)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}

}

func TestAccount_ChangeAccountStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	mockedUsecase := controllers_test.NewMockIAccountUsecase(ctrl)

	validuuid := uuid.NewString()

	testCases := []struct {
		name          string
		in            *models.ChangeAccountStatus
		setupMocks    func()
		out           *models.ChangeAccountStatusResponse
		expectedError error
	}{
		{
			name: "Valid case",
			in: &models.ChangeAccountStatus{
				UUID:   validuuid,
				Status: 1,
			},
			setupMocks: func() {
				mockedUsecase.EXPECT().ChangeAccountStatus(ctx, gomock.Any(), gomock.Any()).Return(nil)
			},
			out: &models.ChangeAccountStatusResponse{
				UUID:   validuuid,
				Status: 1,
			},
			expectedError: nil,
		},
		{
			name: "Invalid UUID",
			in: &models.ChangeAccountStatus{
				UUID:   "invaliduuid",
				Status: 1,
			},
			setupMocks: func() {
			},
			out:           nil,
			expectedError: validators.ErrUUIDIsNotValid,
		},
		{
			name: "Invalid Status",
			in: &models.ChangeAccountStatus{
				UUID:   validuuid,
				Status: 11,
			},
			setupMocks: func() {
			},
			out:           nil,
			expectedError: validators.ErrStatusIsNotValid,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			account := &Account{
				usecase: mockedUsecase,
			}

			result, err := account.ChangeAccountStatus(ctx, tc.in)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &models.ChangeAccountStatusResponse{}, result)
				assert.Equal(t, tc.out.UUID, result.UUID)
				assert.Equal(t, tc.out.Status, result.Status)
				assert.NotEqual(t, 0, result.UpdatedAt)
			}
		})
	}
}
