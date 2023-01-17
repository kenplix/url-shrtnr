package service

import (
	"context"

	"github.com/kenplix/url-shrtnr/pkg/hash"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/repository"
)

type usersService struct {
	usersRepo  repository.UsersRepository
	hasherServ hash.HasherService
}

func NewUsersService(
	usersRepo repository.UsersRepository,
	hasherServ hash.HasherService,
) (UsersService, error) {
	if usersRepo == nil {
		return nil, errors.New("users repository not provided")
	}

	if hasherServ == nil {
		return nil, errors.New("hasher service not provided")
	}

	s := &usersService{
		usersRepo:  usersRepo,
		hasherServ: hasherServ,
	}

	return s, nil
}

func (s *usersService) GetByID(ctx context.Context, userID primitive.ObjectID) (entity.User, error) {
	user, err := s.usersRepo.FindByID(ctx, userID)
	if err != nil {
		return entity.User{}, errors.Wrapf(err, "failed to get user[id:%q]", userID.Hex())
	}

	return user.Filter(), nil
}

func (s *usersService) ChangeEmail(ctx context.Context, schema ChangeEmailSchema) error {
	err := s.usersRepo.ChangeEmail(ctx, repository.ChangeEmailSchema{
		UserID:   schema.UserID,
		NewEmail: schema.NewEmail,
	})
	if err != nil {
		return errors.Wrapf(err, "user[id:%q]: failed to change email", schema.UserID.Hex())
	}

	return nil
}

func (s *usersService) ChangePassword(ctx context.Context, schema ChangePasswordSchema) error {
	user, err := s.usersRepo.FindByID(ctx, schema.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to get user[id:%q]", user.ID.Hex())
	}

	if ok := s.hasherServ.VerifyPassword(schema.CurrentPassword, user.PasswordHash); !ok {
		return entity.ErrIncorrectCredentials
	}

	passwordHash, err := s.hasherServ.HashPassword(schema.NewPassword)
	if err != nil {
		return errors.Wrapf(err, "failed to hash %q password", schema.NewPassword)
	}

	err = s.usersRepo.ChangePassword(ctx, repository.ChangePasswordSchema{
		UserID:          schema.UserID,
		NewPasswordHash: passwordHash,
	})
	if err != nil {
		return errors.Wrapf(err, "user[id:%q]: failed to change password", schema.UserID)
	}

	return nil
}
