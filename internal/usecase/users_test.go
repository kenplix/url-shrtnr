package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	repoMocks "github.com/Kenplix/url-shrtnr/internal/repository/mocks"
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

	type mockBehavior func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher)

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
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
				hasher.
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
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
				hasher.
					On("HashPassword", mock.Anything).
					Return("<password hash>", nil)

				usersRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "correct work",
			args: args{
				input: testUserSignUpInput(t),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
				hasher.
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
			hasher := hashMocks.NewHasher(t)
			usersServ := NewUsersService(usersRepo, hasher)
			tc.mockBehavior(usersRepo, hasher)

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
		tokens Tokens
		hasErr bool
	}

	type mockBehavior func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher)

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
				tokens: Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
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
				tokens: Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.User{}, assert.AnError)
			},
		},
		{
			name: "wrong password check",
			args: args{
				input: func(t *testing.T) UserSignInInput {
					t.Helper()

					input := testUserSignInInput(t)
					input.Password = "PoorestMan"

					return input
				}(t),
			},
			ret: ret{
				tokens: Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasher.
					On("CheckPasswordHash", mock.Anything, mock.Anything).
					Return(false)
			},
		},
		{
			name: "correct work",
			args: args{
				input: testUserSignInInput(t),
			},
			ret: ret{
				tokens: Tokens{
					AccessToken:  "<access token>",
					RefreshToken: "<refresh token>",
				},
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasher *hashMocks.Hasher) {
				usersRepo.
					On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(entity.User{}, nil)

				hasher.
					On("CheckPasswordHash", mock.Anything, mock.Anything).
					Return(true)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usersRepo := repoMocks.NewUsersRepository(t)
			hasher := hashMocks.NewHasher(t)
			usersServ := NewUsersService(usersRepo, hasher)
			tc.mockBehavior(usersRepo, hasher)

			tokens, err := usersServ.SignIn(context.Background(), tc.args.input)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
				return
			}

			assert.Equal(t, tc.ret.tokens, tokens)
		})
	}
}
