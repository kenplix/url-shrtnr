package service_test

import (
	"context"
	"github.com/Kenplix/url-shrtnr/internal/service"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"testing"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	repoMocks "github.com/Kenplix/url-shrtnr/internal/repository/mocks"
	authMocks "github.com/Kenplix/url-shrtnr/pkg/auth/mocks"
	hashMocks "github.com/Kenplix/url-shrtnr/pkg/hash/mocks"
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

	type mockBehavior func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService)

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
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
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
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
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
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
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
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
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

					input := testUserSignUpSchema(t)
					input.Password = "<unhashable password>"

					return input
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
			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			tokensServ := authMocks.NewTokensService(t)

			authServ, err := service.NewAuthService(usersRepo, hasherServ, tokensServ)
			if err != nil {
				t.Fatalf("failed to create auth service: %s", err)
				return
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
		input service.UserSignInSchema
	}

	type ret struct {
		tokens auth.Tokens
		hasErr bool
	}

	type mockBehavior func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService)

	testUserSignInSchema := func(t *testing.T) service.UserSignInSchema {
		t.Helper()

		return service.UserSignInSchema{
			Login:    "bitcoincreator@gmail.com",
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
			name: "user with such credentials not found",
			args: args{
				input: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)
			},
		},
		{
			name: "failed to find user by login",
			args: args{
				input: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "incorrect password",
			args: args{
				input: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
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
				input: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				tokensServ.
					On("CreateTokens", mock.Anything).
					Return(auth.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				input: testUserSignInSchema(t),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				usersRepo.
					On("FindByLogin", mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				tokensServ.
					On("CreateTokens", mock.Anything).
					Return(auth.Tokens{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			tokensServ := authMocks.NewTokensService(t)

			authServ, err := service.NewAuthService(usersRepo, hasherServ, tokensServ)
			if err != nil {
				t.Fatalf("failed to create auth service: %s", err)
				return
			}

			tc.mockBehavior(usersRepo, hasherServ, tokensServ)

			tokens, err := authServ.SignIn(context.Background(), tc.args.input)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokens, tokens)
		})
	}
}
