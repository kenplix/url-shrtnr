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
	"github.com/Kenplix/url-shrtnr/internal/usecase"
	"github.com/Kenplix/url-shrtnr/internal/usecase/mocks"
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

	type mockBehavior func(usersService *mocks.UsersService)

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
			name: "wrong input data",
			args: args{
				inputBody: `¯\_(ツ)_/¯`,
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: errInvalidInputBody.Error(),
				}),
			},
			mockBehavior: func(usersService *mocks.UsersService) {},
		},
		{
			name: "user already exists",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpInput(t)),
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: entity.ErrUserAlreadyExists.Error(),
				}),
			},
			mockBehavior: func(usersService *mocks.UsersService) {
				usersService.
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
			mockBehavior: func(usersService *mocks.UsersService) {
				usersService.
					On("SignUp", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "correct work",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpInput(t)),
			},
			ret: ret{
				statusCode:   http.StatusCreated,
				responseBody: "",
			},
			mockBehavior: func(usersService *mocks.UsersService) {
				usersService.
					On("SignUp", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := mocks.NewUsersService(t)
			handler := NewUsersHandler(service)
			tc.mockBehavior(service)

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

	type mockBehavior func(usersService *mocks.UsersService)

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
			name: "wrong input data",
			args: args{
				inputBody: `¯\_(ツ)_/¯`,
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: errInvalidInputBody.Error(),
				}),
			},
			mockBehavior: func(usersService *mocks.UsersService) {},
		},
		{
			name: "user doesn't exists",
			args: args{
				inputBody: mustMarshal(t, testUserSignInInput(t)),
			},
			ret: ret{
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, ErrorResponse{
					Message: entity.ErrIncorrectEmailOrPassword.Error(),
				}),
			},
			mockBehavior: func(usersService *mocks.UsersService) {
				usersService.
					On("SignIn", mock.Anything, mock.Anything).
					Return(usecase.Tokens{}, entity.ErrIncorrectEmailOrPassword)
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
			mockBehavior: func(usersService *mocks.UsersService) {
				usersService.
					On("SignIn", mock.Anything, mock.Anything).
					Return(usecase.Tokens{}, assert.AnError)
			},
		},
		{
			name: "correct work",
			args: args{
				inputBody: mustMarshal(t, testUserSignInInput(t)),
			},
			ret: ret{
				statusCode: http.StatusOK,
				responseBody: mustMarshal(t, tokenResponse{
					AccessToken:  "<access token>",
					RefreshToken: "<refresh token>",
				}),
			},
			mockBehavior: func(usersService *mocks.UsersService) {
				usersService.
					On("SignIn", mock.Anything, mock.Anything).
					Return(
						usecase.Tokens{
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
			service := mocks.NewUsersService(t)
			handler := NewUsersHandler(service)
			tc.mockBehavior(service)

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

func mustMarshal(t *testing.T, data interface{}) string {
	t.Helper()

	buf, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(buf)
}
