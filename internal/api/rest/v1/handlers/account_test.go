package handlers

import (
	"bytes"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	resthandlers_test "github.com/alexsibrin/runbot-auth/internal/api/rest/v1/handlers/mocks"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedController := resthandlers_test.NewMockIAccountController(ctrl)

	type testCase struct {
		name           string
		in             string
		setupMocks     func()
		expectedBody   string
		expectedCookie string
		expectedCode   int
	}

	testCases := []testCase{
		{
			name: "Valid sign in",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&models.SignInResponse{
					Account: &models.Account{
						UUID:      "someuuid",
						Email:     "some@correctemail.com",
						Password:  "somestrongpswd",
						Name:      "HelloName",
						CreatedAt: time.Now().Unix(),
						UpdatedAt: 0,
					},
					Token: &models.Token{
						Access:  "accesstoken",
						Refresh: "refreshtoken",
					},
				}, nil)
			},
			expectedBody:   `"UUID":"someuuid","Email":"some@correctemail.com","Password":"somestrongpswd","Name":"HelloName"`,
			expectedCookie: `rt=refreshtoken`,
			expectedCode:   200,
		},
		{
			name: "Wrong password",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrPasswordIsWrong)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   400,
		},
		{
			name: "Wrong email",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrEmailIsWrong)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupMocks()

			handler, err := NewAccount(&DependenciesAccount{
				AccountController: mockedController,
				Logger:            logrus.New(),
			})
			assert.NoError(t, err)
			assert.IsType(t, &Account{}, handler)

			req, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer([]byte(tc.in)))
			w := httptest.NewRecorder()

			gctx, _ := gin.CreateTestContext(w)

			gctx.Request = req

			handler.SignIn(gctx)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)

			if tc.expectedCookie == "" {
				return
			}

			var cookiefound bool

			for _, cookie := range w.Result().Cookies() {
				cookiefound = assert.Contains(t, tc.expectedCookie, cookie.Value) && assert.Contains(t, tc.expectedCookie, cookie.Name)
			}

			if !cookiefound {
				t.Error("refresh token is not found")
			}

		})
	}
}

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedController := resthandlers_test.NewMockIAccountController(ctrl)

	type testCase struct {
		name           string
		in             string
		setupMocks     func()
		expectedBody   string
		expectedCookie string
		expectedCode   int
	}

	testCases := []testCase{
		{
			name: "Valid sign up",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd","Name":"SomeName"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&models.SignUpResponse{
					Account: &models.Account{
						UUID:      "someuuid",
						Email:     "some@correctemail.com",
						Password:  "somestrongpswd",
						Name:      "HelloName",
						CreatedAt: time.Now().Unix(),
						UpdatedAt: 0,
					},
					Token: &models.Token{
						Access:  "accesstoken",
						Refresh: "refreshtoken",
					},
				}, nil)
			},
			expectedBody:   `"UUID":"someuuid","Email":"some@correctemail.com","Password":"somestrongpswd","Name":"HelloName"`,
			expectedCookie: `rt=refreshtoken`,
			expectedCode:   200,
		},
		{
			name: "Wrong password",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrPasswordIsWrong)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   400,
		},
		{
			name: "Wrong email",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrEmailIsWrong)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   400,
		},
		{
			name: "Wrong name",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, validators.ErrNameFormatIsNotCorrect)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   400,
		},
		{
			name: "Account is already exist",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			setupMocks: func() {
				mockedController.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrAccountAlreadyExist)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupMocks()

			handler, err := NewAccount(&DependenciesAccount{
				AccountController: mockedController,
				Logger:            logrus.New(),
			})
			assert.NoError(t, err)
			assert.IsType(t, &Account{}, handler)

			req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(tc.in)))
			w := httptest.NewRecorder()

			gctx, _ := gin.CreateTestContext(w)

			gctx.Request = req

			handler.SignUp(gctx)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)

			if tc.expectedCookie == "" {
				return
			}

			var cookiefound bool

			for _, cookie := range w.Result().Cookies() {
				cookiefound = assert.Contains(t, tc.expectedCookie, cookie.Value) && assert.Contains(t, tc.expectedCookie, cookie.Name)
			}

			if !cookiefound {
				t.Error("refresh token is not found")
			}

		})
	}
}
