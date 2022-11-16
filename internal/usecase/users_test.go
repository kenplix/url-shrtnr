package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	repoMocks "github.com/Kenplix/url-shrtnr/internal/repository/mocks"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	authMocks "github.com/Kenplix/url-shrtnr/pkg/auth/mocks"
	hashMocks "github.com/Kenplix/url-shrtnr/pkg/hash/mocks"
)

func TestUsersService_SignUp(t *testing.T) {
	t.Parallel()

	type args struct {
		input UserSignUpInput
	}

	type ret struct {
		hasErr bool
	}

	type mockBehavior func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService)

	testUserSignUpInput := func(t *testing.T) UserSignUpInput {
		t.Helper()

		return UserSignUpInput{
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
			name: "password hashing error",
			args: args{
				input: func(t *testing.T) UserSignUpInput {
					t.Helper()

					input := testUserSignUpInput(t)
					input.Password = "try to imagine a unhashable password"

					return input
				}(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				hasherServ.
					On("HashPassword", mock.Anything).
					Return("", assert.AnError)
			},
		},
		{
			name: "user creation error",
			args: args{
				input: testUserSignUpInput(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
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
				input: testUserSignUpInput(t),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				hasherServ.
					On("HashPassword", mock.Anything).
					Return("<password hash>", nil)

				usersRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			tokensServ := authMocks.NewTokensService(t)

			usersServ := NewUsersService(usersRepo, hasherServ, tokensServ)
			tc.mockBehavior(usersRepo, hasherServ)

			err := usersServ.SignUp(context.Background(), tc.args.input)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}
		})
	}
}

func TestUsersService_SignIn(t *testing.T) {
	t.Parallel()

	type args struct {
		input UserSignInInput
	}

	type ret struct {
		tokens auth.Tokens
		hasErr bool
	}

	type mockBehavior func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService)

	testUserSignInInput := func(t *testing.T) UserSignInInput {
		t.Helper()

		return UserSignInInput{
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
			name: "user not found",
			args: args{
				input: testUserSignInInput(t),
			},
			ret: ret{
				tokens: auth.Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.User{}, entity.ErrUserNotFound)
			},
		},
		{
			name: "getting user error",
			args: args{
				input: testUserSignInInput(t),
			},
			ret: ret{
				tokens: auth.Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "different passwords",
			args: args{
				input: func(t *testing.T) UserSignInInput {
					t.Helper()

					input := testUserSignInInput(t)
					input.Password = "PoorestMan"

					return input
				}(t),
			},
			ret: ret{
				tokens: auth.Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				const userPasswordHash = "<user password hash>"

				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.User{PasswordHash: userPasswordHash}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, userPasswordHash).
					Return(false)
			},
		},
		{
			name: "tokens creation error",
			args: args{
				input: testUserSignInInput(t),
			},
			ret: ret{
				tokens: auth.Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				userID := primitive.NewObjectID()
				const userPasswordHash = "<user password hash>"

				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(
						entity.User{
							ID:           userID,
							PasswordHash: userPasswordHash,
						},
						nil,
					)

				hasherServ.
					On("VerifyPassword", mock.Anything, userPasswordHash).
					Return(true)

				tokensServ.
					On("CreateTokens", userID.Hex()).
					Return(auth.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				input: testUserSignInInput(t),
			},
			ret: ret{
				tokens: auth.Tokens{
					AccessToken:  "<access token>",
					RefreshToken: "<refresh token>",
				},
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService, tokensServ *authMocks.TokensService) {
				userID := primitive.NewObjectID()
				const userPasswordHash = "<user password hash>"

				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(
						entity.User{
							ID:           userID,
							PasswordHash: userPasswordHash,
						},
						nil,
					)

				hasherServ.
					On("VerifyPassword", mock.Anything, userPasswordHash).
					Return(true)

				tokensServ.
					On("CreateTokens", userID.Hex()).
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
			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			tokensServ := authMocks.NewTokensService(t)

			usersServ := NewUsersService(usersRepo, hasherServ, tokensServ)
			tc.mockBehavior(usersRepo, hasherServ, tokensServ)

			tokens, err := usersServ.SignIn(context.Background(), tc.args.input)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokens, tokens)
		})
	}
}

func TestUsersService_RefreshTokens(t *testing.T) {
	t.Parallel()

	type args struct {
		refreshToken string
	}

	type ret struct {
		tokens auth.Tokens
		hasErr bool
	}

	type mockBehavior func(tokensServ *authMocks.TokensService)

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "refresh token parsing error",
			args: args{
				refreshToken: "<refresh token>",
			},
			ret: ret{
				tokens: auth.Tokens{},
				hasErr: true,
			},
			mockBehavior: func(tokensServ *authMocks.TokensService) {
				tokensServ.
					On("ParseRefreshToken", mock.Anything).
					Return("", assert.AnError)
			},
		},
		{
			name: "tokens creation error",
			args: args{
				refreshToken: "<refresh token>",
			},
			ret: ret{
				tokens: auth.Tokens{},
				hasErr: true,
			},
			mockBehavior: func(tokensServ *authMocks.TokensService) {
				userID := primitive.NewObjectID().Hex()

				tokensServ.
					On("ParseRefreshToken", mock.Anything).
					Return(userID, nil)

				tokensServ.
					On("CreateTokens", userID).
					Return(auth.Tokens{}, assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				refreshToken: "<refresh token>",
			},
			ret: ret{
				tokens: auth.Tokens{
					AccessToken:  "<new access token>",
					RefreshToken: "<new refresh token>",
				},
				hasErr: false,
			},
			mockBehavior: func(tokensServ *authMocks.TokensService) {
				userID := primitive.NewObjectID().Hex()

				tokensServ.
					On("ParseRefreshToken", mock.Anything).
					Return(userID, nil)

				tokensServ.
					On("CreateTokens", userID).
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
			usersRepo := repoMocks.NewUsersRepository(t)
			hasherServ := hashMocks.NewHasherService(t)
			tokensServ := authMocks.NewTokensService(t)

			usersServ := NewUsersService(usersRepo, hasherServ, tokensServ)
			tc.mockBehavior(tokensServ)

			tokens, err := usersServ.RefreshTokens(context.Background(), tc.args.refreshToken)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokens, tokens)
		})
	}
}
