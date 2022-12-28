package service_test

import (
	"context"
	"github.com/kenplix/url-shrtnr/internal/entity"
	repoMocks "github.com/kenplix/url-shrtnr/internal/repository/mocks"
	"github.com/kenplix/url-shrtnr/internal/service"
	hashMocks "github.com/kenplix/url-shrtnr/pkg/hash/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUsersService_ChangePassword(t *testing.T) {
	type args struct {
		schema service.ChangePasswordSchema
	}

	type ret struct {
		hasErr bool
	}

	type mockBehavior func(*repoMocks.UsersRepository, *hashMocks.HasherService)

	testChangePasswordSchema := func(t *testing.T) service.ChangePasswordSchema {
		t.Helper()

		return service.ChangePasswordSchema{
			UserID:          primitive.NewObjectID(),
			CurrentPassword: "1wE$Rty2",
			NewPassword:     "2ytR$Ew1",
		}
	}

	testCases := []struct {
		name         string
		args         args
		ret          ret
		mockBehavior mockBehavior
	}{
		{
			name: "failed to get user",
			args: args{
				schema: testChangePasswordSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByID", mock.Anything, mock.Anything).
					Return(entity.UserModel{}, assert.AnError)
			},
		},
		{
			name: "incorrect password",
			args: args{
				schema: testChangePasswordSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByID", mock.Anything, mock.Anything).
					Return(entity.UserModel{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(false)
			},
		},
		{
			name: "failed to hash password",
			args: args{
				schema: func(t *testing.T) service.ChangePasswordSchema {
					t.Helper()

					schema := testChangePasswordSchema(t)
					schema.NewPassword = "<unhashable password>"

					return schema
				}(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByID", mock.Anything, mock.Anything).
					Return(entity.UserModel{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				hasherServ.
					On("HashPassword", mock.Anything).
					Return("", assert.AnError)
			},
		},
		{
			name: "failed to change password",
			args: args{
				schema: testChangePasswordSchema(t),
			},
			ret: ret{
				hasErr: true,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByID", mock.Anything, mock.Anything).
					Return(entity.UserModel{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				hasherServ.
					On("HashPassword", mock.Anything).
					Return("<password hash>", nil)

				usersRepo.
					On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
		{
			name: "ok",
			args: args{
				schema: testChangePasswordSchema(t),
			},
			ret: ret{
				hasErr: false,
			},
			mockBehavior: func(usersRepo *repoMocks.UsersRepository, hasherServ *hashMocks.HasherService) {
				usersRepo.
					On("FindByID", mock.Anything, mock.Anything).
					Return(entity.UserModel{}, nil)

				hasherServ.
					On("VerifyPassword", mock.Anything, mock.Anything).
					Return(true)

				hasherServ.
					On("HashPassword", mock.Anything).
					Return("<password hash>", nil)

				usersRepo.
					On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				usersRepo  = repoMocks.NewUsersRepository(t)
				hasherServ = hashMocks.NewHasherService(t)
			)

			usersServ, err := service.NewUsersService(usersRepo, hasherServ)
			require.NoErrorf(t, err, "failed to create users service: %s", err)

			tc.mockBehavior(usersRepo, hasherServ)

			err = usersServ.ChangePassword(context.Background(), tc.args.schema)
			assert.Falsef(t, (err != nil) != tc.ret.hasErr, "expected error: %t, but got: %v", tc.ret.hasErr, err)
		})
	}
}
