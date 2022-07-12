package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/repository/mocks"
)

func TestUsersService_SignUp(t *testing.T) {
	t.Parallel()

	type args struct {
		input UserSignUpInput
	}

	type ret struct {
		hasErr bool
	}

	type mockBehavior func(usersRepository *mocks.UsersRepository)

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "invalid user",
			args: args{
				input: UserSignUpInput{
					FirstName: "Satoshi",
					LastName:  "Nakamoto",
					Email:     "wrong email",
					Password:  "RichestMan",
				},
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepository *mocks.UsersRepository) {
				usersRepository.
					On(
						"Create",
						mock.Anything,
						mock.Anything,
					).
					Return(assert.AnError)
			},
		},
		{
			name: "valid user",
			args: args{
				input: UserSignUpInput{
					FirstName: "Oleksandr",
					LastName:  "Tolstoi",
					Email:     "no-reply@gmail.com",
					Password:  "12345678",
				},
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(usersRepository *mocks.UsersRepository) {
				usersRepository.
					On(
						"Create",
						mock.Anything,
						mock.Anything,
					).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := mocks.NewUsersRepository(t)
			service := NewUsersService(repository)
			tc.mockBehavior(repository)

			err := service.SignUp(context.Background(), tc.args.input)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
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

	type mockBehavior func(usersRepository *mocks.UsersRepository)

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "invalid user",
			args: args{
				input: UserSignInInput{
					Email:    "wrong email",
					Password: "RichestMan",
				},
			},
			ret: ret{
				tokens: Tokens{},
				hasErr: true,
			},
			mockBehavior: func(usersRepository *mocks.UsersRepository) {
				usersRepository.
					On(
						"GetByCredentials",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).
					Return(
						entity.User{},
						assert.AnError,
					)
			},
		},
		{
			name: "valid user",
			args: args{
				input: UserSignInInput{
					Email:    "no-reply@gmail.com",
					Password: "12345678",
				},
			},
			ret: ret{
				tokens: Tokens{
					AccessToken:  "<access token>",
					RefreshToken: "<refresh token>",
				},
				hasErr: false,
			},
			mockBehavior: func(usersRepository *mocks.UsersRepository) {
				usersRepository.
					On(
						"GetByCredentials",
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).
					Return(
						entity.User{},
						nil,
					)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository := mocks.NewUsersRepository(t)
			service := NewUsersService(repository)
			tc.mockBehavior(repository)

			tokens, err := service.SignIn(context.Background(), tc.args.input)
			if (err != nil) != tc.ret.hasErr {
				t.Errorf("expected error: %t, but got: %v.", tc.ret.hasErr, err)
			}

			assert.Equal(t, tc.ret.tokens, tokens)
		})
	}
}
