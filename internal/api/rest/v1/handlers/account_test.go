package handlers

import (
	"bytes"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	resthandlers_test "github.com/alexsibrin/runbot-auth/internal/api/rest/v1/handlers/mocks"
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

type mocks struct {
	c *resthandlers_test.MockIAccountController
}

func TestSignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name           string
		in             string
		mock           func(*mocks)
		expectedBody   string
		expectedCookie string
		expectedCode   int
	}

	testCases := []testCase{
		{
			name: "Valid sign in",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			mock: func(m *mocks) {
				m.c.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&models.SignInResponse{
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
			mock: func(m *mocks) {
				m.c.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrPasswordIsWrong)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   500,
		},
		{
			name: "Wrong email",
			in:   `{"Email":"some@correctemail.com","Password":"somestrongpswd"}`,
			mock: func(m *mocks) {
				m.c.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, usecases.ErrEmailIsWrong)
			},
			expectedBody:   `{"error":"input data is wrong"}`,
			expectedCookie: ``,
			expectedCode:   500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mc := resthandlers_test.NewMockIAccountController(ctrl)

			m := &mocks{
				c: mc,
			}

			tc.mock(m)

			handler, err := NewAccount(&DependenciesAccount{
				AccountController: mc,
				Logger:            logrus.New(),
			})
			assert.NoError(t, err)
			assert.IsType(t, &Account{}, handler)

			req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer([]byte(tc.in)))
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
