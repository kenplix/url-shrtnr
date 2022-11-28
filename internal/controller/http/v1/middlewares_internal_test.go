package v1

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	servMocks "github.com/Kenplix/url-shrtnr/internal/service/mocks"
	"github.com/Kenplix/url-shrtnr/pkg/token"
)

func TestTranslatorMiddleware(t *testing.T) {
	type args struct {
		localeParameter      []string
		acceptLanguageHeader []string
	}

	type ret struct {
		chosenLocale string
	}

	testCases := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "default locale",
			args: args{
				localeParameter:      []string{"ua"},
				acceptLanguageHeader: []string{},
			},
			ret: ret{
				chosenLocale: "en",
			},
		},
		{
			name: "different locales",
			args: args{
				localeParameter:      []string{},
				acceptLanguageHeader: []string{"ru", "en", "ca"},
			},
			ret: ret{
				chosenLocale: "ru",
			},
		},
		{
			name: "priority choice",
			args: args{
				localeParameter:      []string{"ru"},
				acceptLanguageHeader: []string{"en"},
			},
			ret: ret{
				chosenLocale: "ru",
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := testGinContext(rec)

			query := url.Values{"locale": tc.args.localeParameter}
			c.Request.URL.RawQuery = query.Encode()

			for _, acceptLanguage := range tc.args.acceptLanguageHeader {
				c.Request.Header.Add("Accept-Language", acceptLanguage)
			}

			translatorMiddleware(c)

			translator := c.MustGet(translatorContext).(ut.Translator)
			assert.Equal(t, tc.ret.chosenLocale, translator.Locale())
		})
	}
}

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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersServ := servMocks.NewUsersService(t)
			jwtServ := servMocks.NewJWTService(t)

			handler := userIdentityMiddleware(usersServ, jwtServ)

			tc.mockBehavior(usersServ, jwtServ)

			r := gin.New()
			r.POST("/protected", handler)

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

// testGinContext returns gin context mock
func testGinContext(w *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return c
}

func testUnauthorizedErrorResponse(t *testing.T) string {
	return mustMarshal(t, errResponse{
		Errors: []apiError{newUnauthorizedError()},
	})
}

func testSuspendedErrorResponse(t *testing.T) string {
	return mustMarshal(t, errResponse{
		Errors: []apiError{newSuspendedError()},
	})
}
