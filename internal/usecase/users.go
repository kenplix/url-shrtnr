package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/repository"
)

type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) *usersService {
	return &usersService{
		repo: repo,
	}
}

func (s *usersService) SignUp(ctx context.Context, input UserSignUpInput) error {
	err := s.repo.Create(ctx, entity.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PasswordHash: input.Password,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "could not sign up")
	}

	return nil
}

func (s *usersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	_, err := s.repo.GetByCredentials(ctx, input.Email, input.Password)
	if err != nil {
		return Tokens{}, errors.Wrap(err, "could not sign in")
	}

	return Tokens{
		AccessToken:  "<access token>",
		RefreshToken: "<refresh token>",
	}, nil
}
