package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	usecases_test "github.com/alexsibrin/runbot-auth/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAccountInit(t *testing.T) {

	ctrl := gomock.NewController(t)
	repomock := usecases_test.NewMockIAccountRepo(ctrl)
	hashmock := usecases_test.NewMockIPasswordHasher(ctrl)

	testCases := []struct {
		name        string
		in          *AccountDependencies
		outAccount  *Account
		expectedErr error
	}{
		{
			name: "Regular valid case",
			in: &AccountDependencies{
				Repo:           repomock,
				PasswordHasher: hashmock,
			},
			outAccount: &Account{
				repo:           repomock,
				passwordhasher: hashmock,
			},
			expectedErr: nil,
		},
		{
			name: "Repo is nil case",
			in: &AccountDependencies{
				Repo:           nil,
				PasswordHasher: hashmock,
			},
			outAccount:  nil,
			expectedErr: ErrAccountRepoIsNil,
		},
		{
			name: "PasswordHasher is nil case",
			in: &AccountDependencies{
				Repo:           repomock,
				PasswordHasher: nil,
			},
			outAccount:  nil,
			expectedErr: ErrPaswordHasherIsNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, err := NewAccount(tc.in)
			assert.EqualValues(t, uc, tc.outAccount)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecases_test.NewMockIAccountRepo(ctrl)
	mockHasher := usecases_test.NewMockIPasswordHasher(ctrl)

	ctx := context.TODO()
	testAccount := &entities.Account{Email: "test@example.com", Password: "hashedpassword"}

	tests := []struct {
		name        string
		email       string
		pswd        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name:  "Successful Sign-In",
			email: "test@example.com",
			pswd:  "correctpassword",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "test@example.com").Return(testAccount, nil)
				mockHasher.EXPECT().Compare("correctpassword", "hashedpassword").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "Incorrect Email",
			email: "wrong@example.com",
			pswd:  "password",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "wrong@example.com").Return(nil, errors.New("not found"))
			},
			expectedErr: errors.New("not found"),
		},
		{
			name:  "Incorrect Password",
			email: "test@example.com",
			pswd:  "wrongpassword",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "test@example.com").Return(testAccount, nil)
				mockHasher.EXPECT().Compare("wrongpassword", "hashedpassword").Return(errors.New("password mismatch"))
			},
			expectedErr: ErrPasswordIsWrong,
		},
		{
			name:  "Database Error",
			email: "test@example.com",
			pswd:  "correctpassword",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "test@example.com").Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
		{
			name:  "Password Hasher Error",
			email: "test@example.com",
			pswd:  "correctpassword",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "test@example.com").Return(testAccount, nil)
				mockHasher.EXPECT().Compare("correctpassword", "hashedpassword").Return(errors.New("hasher error"))
			},
			expectedErr: ErrPasswordIsWrong,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{repo: mockRepo, passwordhasher: mockHasher}
			_, err := account.SignIn(ctx, tc.email, tc.pswd)

			t.Log(err)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetOneByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecases_test.NewMockIAccountRepo(ctrl)

	ctx := context.TODO()
	testAccount := &entities.Account{Email: "test@example.com"}

	tests := []struct {
		name        string
		email       string
		setupMocks  func()
		expectedErr error
	}{
		{
			name:  "Successful Fetch",
			email: "test@example.com",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "test@example.com").Return(testAccount, nil)
			},
			expectedErr: nil,
		},
		{
			name:  "Not Found",
			email: "nonexistent@example.com",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			expectedErr: errors.New("not found"),
		},
		{
			name:  "Database Error",
			email: "test@example.com",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByEmail(ctx, "test@example.com").Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{repo: mockRepo}
			_, err := account.GetOneByEmail(ctx, tc.email)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetOneByUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecases_test.NewMockIAccountRepo(ctrl)

	ctx := context.TODO()
	testAccount := &entities.Account{UUID: "123-uuid"}

	tests := []struct {
		name        string
		uuid        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Successful Fetch",
			uuid: "123-uuid",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByUUID(ctx, "123-uuid").Return(testAccount, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Not Found",
			uuid: "nonexistent-uuid",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByUUID(ctx, "nonexistent-uuid").Return(nil, errors.New("not found"))
			},
			expectedErr: errors.New("not found"),
		},
		{
			name: "Database Error",
			uuid: "123-uuid",
			setupMocks: func() {
				mockRepo.EXPECT().GetOneByUUID(ctx, "123-uuid").Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{repo: mockRepo}
			acc, err := account.GetOneByUUID(ctx, tc.uuid)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.IsType(t, &entities.Account{}, acc)
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecases_test.NewMockIAccountRepo(ctrl)

	ctx := context.TODO()
	testAccount := &entities.Account{
		Email:    "myemail",
		Name:     "My name is",
		Password: "mypassword",
	}
	testReq := &AccountCreateRequest{
		Email:    testAccount.Email,
		Name:     testAccount.Name,
		Password: testAccount.Password,
	}

	tests := []struct {
		name        string
		req         *AccountCreateRequest
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Successful Account Creation",
			req:  testReq,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, gomock.Any()).Return(false, nil)
				mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(testAccount, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Account Already Exists",
			req:  testReq,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, gomock.Any()).Return(true, nil)
			},
			expectedErr: ErrAccountAlreadyExist,
		},
		{
			name: "Error on Checking Existence",
			req:  testReq,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, gomock.Any()).Return(false, errors.New("existence check error"))
			},
			expectedErr: errors.New("existence check error"),
		},
		{
			name: "Error on Creating Account",
			req:  testReq,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, gomock.Any()).Return(false, nil)
				mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("create error"))
			},
			expectedErr: errors.New("create error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{repo: mockRepo}
			acc, err := account.Create(ctx, tc.req)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.IsType(t, &entities.Account{}, acc)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecases_test.NewMockIAccountRepo(ctrl)
	mockHasher := usecases_test.NewMockIPasswordHasher(ctrl)

	ctx := context.TODO()
	testAccount := &entities.Account{Email: "test@example.com", Password: "password"} // Adjust as needed

	tests := []struct {
		name        string
		account     *entities.Account
		setupMocks  func()
		expectedErr error
	}{
		{
			name:    "Successful Sign-Up",
			account: testAccount,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, testAccount).Return(false, nil)
				mockHasher.EXPECT().Hash(testAccount.Password).Return("hashedpassword", nil)
				mockRepo.EXPECT().Create(ctx, testAccount).Return(testAccount, nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Account Already Exists",
			account: testAccount,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, testAccount).Return(true, nil)
			},
			expectedErr: ErrAccountAlreadyExist,
		},
		{
			name:    "Error on Checking Existence",
			account: testAccount,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, testAccount).Return(false, errors.New("existence check error"))
			},
			expectedErr: errors.New("existence check error"),
		},
		{
			name:    "Error on Hashing Password",
			account: testAccount,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, testAccount).Return(false, nil)
				mockHasher.EXPECT().Hash(testAccount.Password).Return("", errors.New("hash error"))
			},
			expectedErr: errors.New("hash error"),
		},
		{
			name:    "Error on Creating Account",
			account: testAccount,
			setupMocks: func() {
				mockRepo.EXPECT().IsExist(ctx, testAccount).Return(false, nil)
				mockHasher.EXPECT().Hash(testAccount.Password).Return("hashedpassword", nil)
				mockRepo.EXPECT().Create(ctx, testAccount).Return(nil, errors.New("create error"))
			},
			expectedErr: errors.New("create error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{repo: mockRepo, passwordhasher: mockHasher}
			acc, err := account.SignUp(ctx, tc.account)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.IsType(t, &entities.Account{}, acc)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccount_ChangeAccountStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := usecases_test.NewMockIAccountRepo(ctrl)
	ctx := context.TODO()

	testCases := []struct {
		name        string
		uuid        string
		status      uint8
		setupMocks  func()
		expectedErr error
	}{
		{
			name:   "Valid case",
			uuid:   "validuuid",
			status: 1,
			setupMocks: func() {
				mockRepo.EXPECT().IsExistByUUID(ctx, gomock.Any()).Return(true, nil)
				mockRepo.EXPECT().SetAccountStatus(ctx, gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Account is not exist",
			uuid:   "validuuid",
			status: 1,
			setupMocks: func() {
				mockRepo.EXPECT().IsExistByUUID(ctx, gomock.Any()).Return(false, nil)
			},
			expectedErr: ErrAccountIsNotExist,
		},
		{
			name:   "Repo error IsExistByUUID",
			uuid:   "validuuid",
			status: 1,
			setupMocks: func() {
				mockRepo.EXPECT().IsExistByUUID(ctx, gomock.Any()).Return(false, fmt.Errorf("some repo error"))
			},
			expectedErr: fmt.Errorf("some repo error"),
		},
		{
			name:   "Repo error SetAccountStatus",
			uuid:   "validuuid",
			status: 1,
			setupMocks: func() {
				mockRepo.EXPECT().IsExistByUUID(ctx, gomock.Any()).Return(true, nil)
				mockRepo.EXPECT().SetAccountStatus(ctx, gomock.Any(), gomock.Any()).Return(fmt.Errorf("some repo error"))
			},
			expectedErr: fmt.Errorf("some repo error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			account := &Account{repo: mockRepo}
			err := account.ChangeAccountStatus(ctx, tc.uuid, tc.status)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
