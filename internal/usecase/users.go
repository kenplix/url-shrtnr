package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
)

type usersService struct {
	usersRepo repository.UsersRepository
	hasher    hash.Hasher
}

func NewUsersService(usersRepo repository.UsersRepository, hasher hash.Hasher) *usersService {
	return &usersService{
		usersRepo: usersRepo,
		hasher:    hasher,
	}
}

func (s *usersService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.HashPassword(input.Password)
	if err != nil {
		return errors.Wrapf(err, "could not hash password %q", input.Password)
	}

	user := entity.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PasswordHash: passwordHash,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
	}

	err = s.usersRepo.Create(ctx, user)
	if err != nil {
		return errors.Wrapf(err, "could not create user %#v", user)
	}

	return nil
}

func (s *usersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	user, err := s.usersRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return Tokens{}, entity.ErrIncorrectEmailOrPassword
		}

		return Tokens{}, errors.Wrapf(err, "could not get user by email %q", input.Email)
	}

	if ok := s.hasher.CheckPasswordHash(input.Password, user.PasswordHash); !ok {
		return Tokens{}, entity.ErrIncorrectEmailOrPassword
	}

	return Tokens{
		AccessToken:  "<access token>",
		RefreshToken: "<refresh token>",
	}, nil
}
