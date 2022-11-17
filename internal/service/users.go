package service

import (
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/pkg/errors"
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
