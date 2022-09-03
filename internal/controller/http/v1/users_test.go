package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	usecaseMocks "github.com/Kenplix/url-shrtnr/internal/usecase/mocks"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	authMocks "github.com/Kenplix/url-shrtnr/pkg/auth/mocks"
)

func TestUsersHandler_UserSignUp(t *testing.T) {
	t.Parallel()

	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(usersServ *usecaseMocks.UsersService)

	testUserSignUpInput := func(t *testing.T) userSignUpInput {
		t.Helper()

		return userSignUpInput{
			FirstName: "Satoshi",
			LastName:  "Nakamoto",
			Email:     "bitcoincreator@gmail.com",
			Password:  "RichestMan",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "invalid input body",
			args: args{
				inputBody: `¯\_(ツ)_/¯`,
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: errInvalidInputBody.Error(),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {},
		},
		{
			name: "user with such email already exists",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpInput(t)),
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: entity.ErrUserAlreadyExists.Error(),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("SignUp", mock.Anything, mock.Anything).
					Return(entity.ErrUserAlreadyExists)
			},
		},
		{
			name: "service failure",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpInput(t)),
			},
			ret: ret{
				statusCode: http.StatusInternalServerError,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("SignUp", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpInput(t)),
			},
			ret: ret{
				statusCode:   http.StatusCreated,
				responseBody: "",
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("SignUp", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersServ := usecaseMocks.NewUsersService(t)
			tokenServ := authMocks.NewTokenService(t)

			handler := NewUsersHandler(usersServ, tokenServ)
			tc.mockBehavior(usersServ)

			r := gin.New()
			r.POST("/sign-up", handler.userSignUp)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(tc.args.inputBody))

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.ret.statusCode, rec.Code)
			assert.Equal(t, tc.ret.responseBody, rec.Body.String())
		})
	}
}

func TestUsersHandler_UserSignIn(t *testing.T) {
	t.Parallel()

	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(usersServ *usecaseMocks.UsersService)

	testUserSignInInput := func(t *testing.T) userSignInInput {
		t.Helper()

		return userSignInInput{
			Email:    "bitcoincreator@gmail.com",
			Password: "RichestMan",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "invalid input body",
			args: args{
				inputBody: `¯\_(ツ)_/¯`,
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: errInvalidInputBody.Error(),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {},
		},
		{
			name: "incorrect email or password",
			args: args{
				inputBody: mustMarshal(t, testUserSignInInput(t)),
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: entity.ErrIncorrectEmailOrPassword.Error(),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(auth.Tokens{}, entity.ErrIncorrectEmailOrPassword)
			},
		},
		{
			name: "service failure",
			args: args{
				inputBody: mustMarshal(t, testUserSignInInput(t)),
			},
			ret: ret{
				statusCode: http.StatusInternalServerError,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(auth.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserSignInInput(t)),
			},
			ret: ret{
				statusCode: http.StatusOK,
				responseBody: mustMarshal(t, auth.Tokens{
					AccessToken:  "<access token>",
					RefreshToken: "<refresh token>",
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(
						auth.Tokens{
							AccessToken:  "<access token>",
							RefreshToken: "<refresh token>",
						},
						nil,
					)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersServ := usecaseMocks.NewUsersService(t)
			tokenServ := authMocks.NewTokenService(t)

			handler := NewUsersHandler(usersServ, tokenServ)
			tc.mockBehavior(usersServ)

			r := gin.New()
			r.POST("/sign-in", handler.userSignIn)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(tc.args.inputBody))

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.ret.statusCode, rec.Code)
			assert.Equal(t, tc.ret.responseBody, rec.Body.String())
		})
	}
}

func TestUsersHandler_UserRefreshTokens(t *testing.T) {
	t.Parallel()

	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(usersServ *usecaseMocks.UsersService)

	testUserRefreshTokensInput := func(t *testing.T) userRefreshTokensInput {
		t.Helper()

		return userRefreshTokensInput{
			RefreshToken: "<refresh token>",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "invalid input body",
			args: args{
				inputBody: `¯\_(ツ)_/¯`,
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: errInvalidInputBody.Error(),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {},
		},
		{
			name: "service failure",
			args: args{
				inputBody: mustMarshal(t, testUserRefreshTokensInput(t)),
			},
			ret: ret{
				statusCode: http.StatusInternalServerError,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("RefreshTokens", mock.Anything, mock.Anything).
					Return(auth.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserRefreshTokensInput(t)),
			},
			ret: ret{
				statusCode: http.StatusOK,
				responseBody: mustMarshal(t, auth.Tokens{
					AccessToken:  "<new access token>",
					RefreshToken: "<new refresh token>",
				}),
			},
			mockBehavior: func(usersServ *usecaseMocks.UsersService) {
				usersServ.
					On("RefreshTokens", mock.Anything, mock.Anything).
					Return(
						auth.Tokens{
							AccessToken:  "<new access token>",
							RefreshToken: "<new refresh token>",
						},
						nil,
					)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersServ := usecaseMocks.NewUsersService(t)
			tokenServ := authMocks.NewTokenService(t)

			handler := NewUsersHandler(usersServ, tokenServ)
			tc.mockBehavior(usersServ)

			r := gin.New()
			r.POST("/refresh-tokens", handler.userRefreshTokens)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/refresh-tokens", bytes.NewBufferString(tc.args.inputBody))

			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.ret.statusCode, rec.Code)
			assert.Equal(t, tc.ret.responseBody, rec.Body.String())
		})
	}
}

func mustMarshal(t *testing.T, data interface{}) string {
	t.Helper()

	buf, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	return string(buf)
}
