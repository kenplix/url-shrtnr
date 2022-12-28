package v1

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/kenplix/url-shrtnr/internal/service"
	servMocks "github.com/kenplix/url-shrtnr/internal/service/mocks"
	"github.com/kenplix/url-shrtnr/pkg/token"
)

func TestHandler_ChangePassword(t *testing.T) {
	type args struct {
		inputBody string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(*servMocks.JWTService, *servMocks.UsersService)

	testUserChangePasswordSchema := func(t *testing.T) userChangePasswordSchema {
		t.Helper()

		return userChangePasswordSchema{
			CurrentPassword:      "1wE$Rty2",
			NewPassword:          "2ytR$Ew1",
			PasswordConfirmation: "2ytR$Ew1",
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
				statusCode: http.StatusBadRequest,
				responseBody: mustMarshal(t, errResponse{
					Errors: []apiError{
						&entity.CoreError{
							Code:    errorcode.InvalidSchema,
							Message: "body should be a JSON object",
						},
					},
				}),
			},
			mockBehavior: func(jwtServ *servMocks.JWTService, usersServ *servMocks.UsersService) {
				userID := primitive.NewObjectID()
				claims := &token.JWTCustomClaims{
					StandardClaims: jwt.StandardClaims{
						Subject: userID.Hex(),
					},
				}

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(claims, nil)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{ID: userID}, nil)

				jwtServ.
					On("ProlongTokens", mock.Anything, mock.Anything)
			},
		},
		{
			name: "incorrect password",
			args: args{
				inputBody: mustMarshal(t, testUserChangePasswordSchema(t)),
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
			mockBehavior: func(jwtServ *servMocks.JWTService, usersServ *servMocks.UsersService) {
				userID := primitive.NewObjectID()
				claims := &token.JWTCustomClaims{
					StandardClaims: jwt.StandardClaims{
						Subject: userID.Hex(),
					},
				}

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(claims, nil)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{ID: userID}, nil)

				jwtServ.
					On("ProlongTokens", mock.Anything, mock.Anything)

				usersServ.
					On("ChangePassword", mock.Anything, mock.Anything).
					Return(entity.ErrIncorrectCredentials)
			},
		},
		{
			name: "service failure",
			args: args{
				inputBody: mustMarshal(t, testUserChangePasswordSchema(t)),
			},
			ret: ret{
				statusCode:   http.StatusInternalServerError,
				responseBody: testInternalErrorResponse(t),
			},
			mockBehavior: func(jwtServ *servMocks.JWTService, usersServ *servMocks.UsersService) {
				userID := primitive.NewObjectID()
				claims := &token.JWTCustomClaims{
					StandardClaims: jwt.StandardClaims{
						Subject: userID.Hex(),
					},
				}

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(claims, nil)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{ID: userID}, nil)

				jwtServ.
					On("ProlongTokens", mock.Anything, mock.Anything)

				usersServ.
					On("ChangePassword", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				inputBody: mustMarshal(t, testUserChangePasswordSchema(t)),
			},
			ret: ret{
				statusCode:   http.StatusOK,
				responseBody: "",
			},
			mockBehavior: func(jwtServ *servMocks.JWTService, usersServ *servMocks.UsersService) {
				userID := primitive.NewObjectID()
				claims := &token.JWTCustomClaims{
					StandardClaims: jwt.StandardClaims{
						Subject: userID.Hex(),
					},
				}

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(claims, nil)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{ID: userID}, nil)

				jwtServ.
					On("ProlongTokens", mock.Anything, mock.Anything)

				usersServ.
					On("ChangePassword", mock.Anything, mock.Anything).
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

			r := gin.New()
			h.InitRoutes(r.Group("/api", testLoggerMiddleware(t)))

			tc.mockBehavior(jwtServ, usersServ)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/users/change-password",
				bytes.NewBufferString(tc.args.inputBody),
			)
			req.Header.Set("Authorization", `Bearer "header.payload.signature"`)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.ret.statusCode, resp.StatusCode)
			assert.Equal(t, tc.ret.responseBody, string(body))
		})
	}
}
