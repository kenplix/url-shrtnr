package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Kenplix/url-shrtnr/internal/service"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	repoMocks "github.com/Kenplix/url-shrtnr/internal/repository/mocks"
	servMocks "github.com/Kenplix/url-shrtnr/internal/service/mocks"
	hashMocks "github.com/Kenplix/url-shrtnr/pkg/hash/mocks"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_SignUp(t *testing.T) {
	type args struct {
		schema service.UserSignUpSchema
	}

	type ret struct {
		hasErr bool
	}

	type mockBehavior func(*repoMocks.UsersRepository, *hashMocks.HasherService)

	testUserSignUpSchema := func(t *testing.T) service.UserSignUpSchema {
		t.Helper()

		return service.UserSignUpSchema{
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
			name: "user with such email already exists",
			args: args{
				schema: testUserSignUpSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, _ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)
			},
		},
		{
			name: "failed to find user by email",
			args: args{
				schema: testUserSignUpSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, _ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "user with such username already exists",
			args: args{
				schema: testUserSignUpSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, _ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				usersRepo.
					On("FindByUsername", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)
			},
		},
		{
			name: "failed to find user by username",
			args: args{
				schema: testUserSignUpSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, _ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				usersRepo.
					On("FindByUsername", mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "failed to hash %q password",
			args: args{
				schema: func(t *testing.T) service.UserSignUpSchema {
					t.Helper()

					schema := testUserSignUpSchema(t)
					schema.Password = "<unhashable password>"

					return schema
				}(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				usersRepo.
					On("FindByUsername", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				hasherServ.
					On("HashPassword", mock.Anything).
					Return("", assert.AnError)
			},
		},
		{
			name: "failed to create user",
			args: args{
				schema: testUserSignUpSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				usersRepo.
					On("FindByUsername", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				hasherServ.
					On("HashPassword", mock.Anything).
					Return("<password hash>", nil)

				usersRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				schema: testUserSignUpSchema(t),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByEmail", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				usersRepo.
					On("FindByUsername", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)

				hasherServ.
					On("HashPassword", mock.Anything).
					Return("<password hash>", nil)

				usersRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redisServ := miniredis.RunT(t)
			cache := redis.NewClient(&redis.Options{
				Addr: redisServ.Addr(),
			})

			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			jwtServ := servMocks.NewJWTService(t)

			authServ, err := service.NewAuthService(cache, usersRepo, hasherServ, jwtServ)
			if err != nil {
				t.Fatalf("failed to create auth service: %s", err)
			}

			tc.mockBehavior(usersRepo, hasherServ)

			err = authServ.SignUp(context.Background(), tc.args.schema)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}
		})
	}
}

func TestAuthService_SignIn(t *testing.T) {
	type args struct {
		schema service.UserSignInSchema
	}

	type ret struct {
		tokens entity.Tokens
		hasErr bool
	}

	type mockBehavior func(
		*repoMocks.UsersRepository,
		*hashMocks.HasherService,
		*servMocks.JWTService,
	)

	testUserSignInSchema := func(t *testing.T) service.UserSignInSchema {
		t.Helper()

		return service.UserSignInSchema{
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
			name: "user with such credentials not found",
			args: args{
				schema: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(
				usersRepo *repoMocks.UsersRepository,
				_ *hashMocks.HasherService,
				_ *servMocks.JWTService,
			) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)
			},
		},
		{
			name: "failed to find user by login",
			args: args{
				schema: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(
				usersRepo *repoMocks.UsersRepository,
				_ *hashMocks.HasherService,
				_ *servMocks.JWTService,
			) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "suspended user",
			args: args{
				schema: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(
				usersRepo *repoMocks.UsersRepository,
				_ *hashMocks.HasherService,
				_ *servMocks.JWTService,
			) {
				suspendedAt := time.Now()

				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{SuspendedAt: &suspendedAt}, nil)
			},
		},
		{
			name: "incorrect password",
			args: args{
				schema: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(
				usersRepo *repoMocks.UsersRepository,
				hasherServ *hashMocks.HasherService,
				_ *servMocks.JWTService,
			) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(false)
			},
		},
		{
			name: "failed to create tokens",
			args: args{
				schema: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(
				usersRepo *repoMocks.UsersRepository,
				hasherServ *hashMocks.HasherService,
				jwtServ *servMocks.JWTService,
			) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				jwtServ.
					On("CreateTokens", mock.Anything, mock.Anything).
					Return(entity.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				schema: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(
				usersRepo *repoMocks.UsersRepository,
				hasherServ *hashMocks.HasherService,
				jwtServ *servMocks.JWTService,
			) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				jwtServ.
					On("CreateTokens", mock.Anything, mock.Anything).
					Return(entity.Tokens{}, nil)
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redisServ := miniredis.RunT(t)
			cache := redis.NewClient(&redis.Options{
				Addr: redisServ.Addr(),
			})

			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			jwtServ := servMocks.NewJWTService(t)

			authServ, err := service.NewAuthService(cache, usersRepo, hasherServ, jwtServ)
			if err != nil {
				t.Fatalf("failed to create auth service: %s", err)
			}

			tc.mockBehavior(usersRepo, hasherServ, jwtServ)

			tokens, err := authServ.SignIn(context.Background(), tc.args.schema)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokens, tokens)
		})
	}
}

func TestAuthService_SignOut(t *testing.T) {
	type args struct {
		userID primitive.ObjectID
	}

	type ret struct {
		hasErr bool
	}

	type mockBehavior func(userID primitive.ObjectID) func(*miniredis.Miniredis)

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "user already signed out",
			args: args{
				userID: primitive.NewObjectID(),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(userID primitive.ObjectID) func(*miniredis.Miniredis) {
				return func(redisServ *miniredis.Miniredis) {}
			},
		},
		{
			name: "ok",
			args: args{
				userID: primitive.NewObjectID(),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(userID primitive.ObjectID) func(*miniredis.Miniredis) {
				return func(redisServ *miniredis.Miniredis) {
					err := redisServ.Set(service.TokenCacheKey(userID.Hex()), mustMarshal(t, entity.TokensUIDs{
						AccessTokenUID:  "<access token UID>",
						RefreshTokenUID: "<refresh token UID>",
					}))
					if err != nil {
						t.Fatalf("failed to set %q token cache key: %s", service.TokenCacheKey(userID.Hex()), err)
					}
				}
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redisServ := miniredis.RunT(t)
			cache := redis.NewClient(&redis.Options{
				Addr: redisServ.Addr(),
			})

			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			jwtServ := servMocks.NewJWTService(t)

			authServ, err := service.NewAuthService(cache, usersRepo, hasherServ, jwtServ)
			if err != nil {
				t.Fatalf("failed to create auth service: %s", err)
			}

			tc.mockBehavior(tc.args.userID)(redisServ)

			err = authServ.SignOut(context.Background(), tc.args.userID)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}
		})
	}
}

func mustMarshal(t *testing.T, data interface{}) string {
	t.Helper()

	buf, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal %v data", err)
	}

	return string(buf)
}
