package v1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/kenplix/url-shrtnr/internal/service"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/kenplix/url-shrtnr/pkg/token"

	"github.com/stretchr/testify/assert"

	servMocks "github.com/kenplix/url-shrtnr/internal/service/mocks"
)

func TestAuthHandler_SignUp(t *testing.T) {
	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(*servMocks.AuthService)

	testUserSignUpSchema := func(t *testing.T) userSignUpSchema {
		t.Helper()

		return userSignUpSchema{
			Username: "kenplix",
			Email:    "tolstoi.job@gmail.com",
			Password: "1wE$Rty2",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "binding error",
			args: args{
				inputBody: "[]",
			},
			ret: ret{
				statusCode:   http.StatusBadRequest,
				responseBody: testUnmarshalTypeError(t),
			},
			mockBehavior: func(_ *servMocks.AuthService) {},
		},
		{
			name: "user already exists",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: mustMarshal(t, errResponse{
					Errors: []apiError{
						&entity.ValidationError{
							CoreError: entity.CoreError{
								Code:    errorcode.AlreadyExists,
								Message: "username is already taken",
							},
							Field: "username",
						},
					},
				}),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignUp", mock.Anything, mock.Anything).
					Return(&entity.ValidationError{
						CoreError: entity.CoreError{
							Code:    errorcode.AlreadyExists,
							Message: "username is already taken",
						},
						Field: "username",
					})
			},
		},
		{
			name: "service failure",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpSchema(t)),
			},
			ret: ret{
				statusCode:   http.StatusInternalServerError,
				responseBody: testInternalErrorResponse(t),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignUp", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserSignUpSchema(t)),
			},
			ret: ret{
				statusCode:   http.StatusCreated,
				responseBody: "",
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignUp", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}

	t.Parallel()

	initValidator(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				jwtServ   = servMocks.NewJWTService(t)
				authServ  = servMocks.NewAuthService(t)
				usersServ = servMocks.NewUsersService(t)
			)

			h, err := NewHandler(testLogger(t), &service.Services{
				JWT:   jwtServ,
				Auth:  authServ,
				Users: usersServ,
			})
			require.NoErrorf(t, err, "failed to create handler: %s", err)

			tc.mockBehavior(authServ)

			r := gin.New()
			r.POST("/sign-up", testLoggerMiddleware(t), h.signUp)

			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(tc.args.inputBody))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.ret.statusCode, resp.StatusCode)
			assert.Equal(t, tc.ret.responseBody, string(body))
		})
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(*servMocks.AuthService)

	testUserSignInSchema := func(t *testing.T) userSignInSchema {
		t.Helper()

		return userSignInSchema{
			Login:    "tolstoi.job@gmail.com",
			Password: "1wE$Rty2",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "binding error",
			args: args{
				inputBody: "[]",
			},
			ret: ret{
				statusCode:   http.StatusBadRequest,
				responseBody: testUnmarshalTypeError(t),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {},
		},
		{
			name: "incorrect email or password",
			args: args{
				inputBody: mustMarshal(t, testUserSignInSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: mustMarshal(t, errResponse{
					Errors: []apiError{
						&entity.CoreError{
							Code:    errorcode.IncorrectCredentials,
							Message: entity.ErrIncorrectCredentials.Error(),
						},
					},
				}),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(entity.Tokens{}, entity.ErrIncorrectCredentials)
			},
		},
		{
			name: "suspended user",
			args: args{
				inputBody: mustMarshal(t, testUserSignInSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusForbidden,
				responseBody: mustMarshal(t, errResponse{
					Errors: []apiError{
						&entity.CoreError{
							Code:    errorcode.CurrentUserSuspended,
							Message: "your account has been suspended",
						},
					},
				}),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(entity.Tokens{}, &entity.SuspendedUserError{UserID: "<user id>"})
			},
		},
		{
			name: "service failure",
			args: args{
				inputBody: mustMarshal(t, testUserSignInSchema(t)),
			},
			ret: ret{
				statusCode:   http.StatusInternalServerError,
				responseBody: testInternalErrorResponse(t),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(entity.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserSignInSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusOK,
				responseBody: mustMarshal(t, entity.Tokens{
					AccessToken:  "<access token>",
					RefreshToken: "<refresh token>",
				}),
			},
			mockBehavior: func(authServ *servMocks.AuthService) {
				authServ.
					On("SignIn", mock.Anything, mock.Anything).
					Return(
						entity.Tokens{
							AccessToken:  "<access token>",
							RefreshToken: "<refresh token>",
						},
						nil,
					)
			},
		},
	}

	t.Parallel()

	initValidator(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				jwtServ   = servMocks.NewJWTService(t)
				authServ  = servMocks.NewAuthService(t)
				usersServ = servMocks.NewUsersService(t)
			)

			h, err := NewHandler(testLogger(t), &service.Services{
				JWT:   jwtServ,
				Auth:  authServ,
				Users: usersServ,
			})
			require.NoErrorf(t, err, "failed to create handler: %s", err)

			tc.mockBehavior(authServ)

			r := gin.New()
			r.POST("/sign-in", testLoggerMiddleware(t), h.signIn)

			req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewBufferString(tc.args.inputBody))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.ret.statusCode, resp.StatusCode)
			assert.Equal(t, tc.ret.responseBody, string(body))
		})
	}
}

func TestAuthHandler_RefreshTokens(t *testing.T) {
	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(*servMocks.JWTService)

	testUserRefreshTokensSchema := func(t *testing.T) userRefreshTokensSchema {
		t.Helper()

		return userRefreshTokensSchema{
			RefreshToken: "header.payload.signature",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "binding error",
			args: args{
				inputBody: "[]",
			},
			ret: ret{
				statusCode:   http.StatusBadRequest,
				responseBody: testUnmarshalTypeError(t),
			},
			mockBehavior: func(_ *servMocks.JWTService) {},
		},
		{
			name: "token parsing error",
			args: args{
				inputBody: mustMarshal(t, testUserRefreshTokensSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: mustMarshal(t, errResponse{
					Errors: []apiError{
						&entity.ValidationError{
							CoreError: entity.CoreError{
								Code:    errorcode.InvalidField,
								Message: "refresh token is invalid, expired or revoked",
							},
							Field: "refreshToken",
						},
					},
				}),
			},
			mockBehavior: func(jwtServ *servMocks.JWTService) {
				jwtServ.
					On("ParseRefreshToken", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
		{
			name: "token validating error",
			args: args{
				inputBody: mustMarshal(t, testUserRefreshTokensSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusUnprocessableEntity,
				responseBody: mustMarshal(t, errResponse{
					Errors: []apiError{
						&entity.ValidationError{
							CoreError: entity.CoreError{
								Code:    errorcode.InvalidField,
								Message: "refresh token is invalid, expired or revoked",
							},
							Field: "refreshToken",
						},
					},
				}),
			},
			mockBehavior: func(jwtServ *servMocks.JWTService) {
				jwtServ.
					On("ParseRefreshToken", mock.Anything).
					Return(&token.JWTCustomClaims{}, nil)

				jwtServ.
					On("ValidateRefreshToken", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "tokens creation error",
			args: args{
				inputBody: mustMarshal(t, testUserRefreshTokensSchema(t)),
			},
			ret: ret{
				statusCode:   http.StatusInternalServerError,
				responseBody: testInternalErrorResponse(t),
			},
			mockBehavior: func(jwtServ *servMocks.JWTService) {
				jwtServ.
					On("ParseRefreshToken", mock.Anything).
					Return(&token.JWTCustomClaims{}, nil)

				jwtServ.
					On("ValidateRefreshToken", mock.Anything, mock.Anything).
					Return(nil)

				jwtServ.
					On("CreateTokens", mock.Anything, mock.Anything).
					Return(entity.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserRefreshTokensSchema(t)),
			},
			ret: ret{
				statusCode: http.StatusOK,
				responseBody: mustMarshal(t, entity.Tokens{
					AccessToken:  "<new access token>",
					RefreshToken: "<new refresh token>",
				}),
			},
			mockBehavior: func(jwtServ *servMocks.JWTService) {
				jwtServ.
					On("ParseRefreshToken", mock.Anything).
					Return(&token.JWTCustomClaims{}, nil)

				jwtServ.
					On("ValidateRefreshToken", mock.Anything, mock.Anything).
					Return(nil)

				jwtServ.
					On("CreateTokens", mock.Anything, mock.Anything).
					Return(
						entity.Tokens{
							AccessToken:  "<new access token>",
							RefreshToken: "<new refresh token>",
						},
						nil,
					)
			},
		},
	}

	t.Parallel()

	initValidator(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				jwtServ   = servMocks.NewJWTService(t)
				authServ  = servMocks.NewAuthService(t)
				usersServ = servMocks.NewUsersService(t)
			)

			h, err := NewHandler(testLogger(t), &service.Services{
				JWT:   jwtServ,
				Auth:  authServ,
				Users: usersServ,
			})
			require.NoErrorf(t, err, "failed to create handler: %s", err)

			tc.mockBehavior(jwtServ)

			r := gin.New()
			r.POST("/refresh-tokens", testLoggerMiddleware(t), h.refreshTokens)

			req := httptest.NewRequest(http.MethodPost, "/refresh-tokens", bytes.NewBufferString(tc.args.inputBody))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.ret.statusCode, resp.StatusCode)
			assert.Equal(t, tc.ret.responseBody, string(body))
		})
	}
}
