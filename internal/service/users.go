package service

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/repository"
)

type usersService struct {
	usersRepo repository.UsersRepository
}

func NewUsersService(usersRepo repository.UsersRepository) (UsersService, error) {
	if usersRepo == nil {
		return nil, errors.New("users repository not provided")
	}

	s := &usersService{
		usersRepo: usersRepo,
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
