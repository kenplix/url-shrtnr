package v1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kenplix/url-shrtnr/internal/service"
	servMocks "github.com/kenplix/url-shrtnr/internal/service/mocks"
)

func TestUserIdentityMiddleware(t *testing.T) {
	type args struct {
		authHeader string
	}

	type ret struct {
		statusCode   int
		responseBody string
	}

	type mockBehavior func(*servMocks.UsersService, *servMocks.JWTService)

	testAuthorizationHeader := func(t *testing.T) string {
		return "Bearer <token>"
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: `empty "Authorization" header`,
			args: args{
				authHeader: "",
			},
			ret: ret{
				statusCode:   http.StatusUnauthorized,
				responseBody: testUnauthorizedErrorResponse(t),
			},
			mockBehavior: func(_ *servMocks.UsersService, _ *servMocks.JWTService) {},
		},
		{
			name: `invalid "Authorization" header`,
			args: args{
				authHeader: "Bearer",
			},
			ret: ret{
				statusCode:   http.StatusUnauthorized,
				responseBody: testUnauthorizedErrorResponse(t),
			},
			mockBehavior: func(_ *servMocks.UsersService, _ *servMocks.JWTService) {},
		},
		{
			name: "failed to parse access token",
			args: args{
				authHeader: testAuthorizationHeader(t),
			},
			ret: ret{
				statusCode:   http.StatusUnauthorized,
				responseBody: testUnauthorizedErrorResponse(t),
			},
			mockBehavior: func(_ *servMocks.UsersService, jwtServ *servMocks.JWTService) {
				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
		{
			name: "failed to validate access token",
			args: args{
				authHeader: testAuthorizationHeader(t),
			},
			ret: ret{
				statusCode:   http.StatusUnauthorized,
				responseBody: testUnauthorizedErrorResponse(t),
			},
			mockBehavior: func(usersServ *servMocks.UsersService, jwtServ *servMocks.JWTService) {
				userID := primitive.NewObjectID()

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(
						&token.JWTCustomClaims{
							StandardClaims: jwt.StandardClaims{
								Subject: userID.Hex(),
							},
						},
						nil,
					)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(assert.AnError)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{ID: userID}, nil)
			},
		},
		{
			name: "failed to parse userID object",
			args: args{
				authHeader: testAuthorizationHeader(t),
			},
			ret: ret{
				statusCode:   http.StatusUnauthorized,
				responseBody: testUnauthorizedErrorResponse(t),
			},
			mockBehavior: func(usersServ *servMocks.UsersService, jwtServ *servMocks.JWTService) {
				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(&token.JWTCustomClaims{}, nil)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
		{
			name: "failed to get user",
			args: args{
				authHeader: testAuthorizationHeader(t),
			},
			ret: ret{
				statusCode:   http.StatusUnauthorized,
				responseBody: testUnauthorizedErrorResponse(t),
			},
			mockBehavior: func(usersServ *servMocks.UsersService, jwtServ *servMocks.JWTService) {
				userID := primitive.NewObjectID()

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(
						&token.JWTCustomClaims{
							StandardClaims: jwt.StandardClaims{
								Subject: userID.Hex(),
							},
						},
						nil,
					)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "suspended user",
			args: args{
				authHeader: testAuthorizationHeader(t),
			},
			ret: ret{
				statusCode:   http.StatusForbidden,
				responseBody: testSuspendedErrorResponse(t),
			},
			mockBehavior: func(usersServ *servMocks.UsersService, jwtServ *servMocks.JWTService) {
				userID := primitive.NewObjectID()

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(
						&token.JWTCustomClaims{
							StandardClaims: jwt.StandardClaims{
								Subject: userID.Hex(),
							},
						},
						nil,
					)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				suspendedAt := time.Now()

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(
						entity.User{
							ID:          userID,
							SuspendedAt: &suspendedAt,
						},
						nil,
					)
			},
		},
		{
			name: "ok",
			args: args{
				authHeader: testAuthorizationHeader(t),
			},
			ret: ret{
				statusCode:   http.StatusOK,
				responseBody: "",
			},
			mockBehavior: func(usersServ *servMocks.UsersService, jwtServ *servMocks.JWTService) {
				userID := primitive.NewObjectID()

				jwtServ.
					On("ParseAccessToken", mock.Anything).
					Return(
						&token.JWTCustomClaims{
							StandardClaims: jwt.StandardClaims{
								Subject: userID.Hex(),
							},
						},
						nil,
					)

				jwtServ.
					On("ValidateAccessToken", mock.Anything, mock.Anything).
					Return(nil)

				usersServ.
					On("GetByID", mock.Anything, mock.Anything).
					Return(entity.User{ID: userID}, nil)
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
			require.NoError(t, err, "failed to create handler: %s", err)

			tc.mockBehavior(usersServ, jwtServ)

			r := gin.New()
			r.POST("/protected", testLoggerMiddleware(t), h.userIdentityMiddleware)

			req := httptest.NewRequest(http.MethodPost, "/protected", bytes.NewBufferString(""))
			req.Header.Set("Authorization", tc.args.authHeader)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tc.ret.statusCode, resp.StatusCode)
			assert.Equal(t, tc.ret.responseBody, string(body))
		})
	}
}
