package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
)

type usersService struct {
	usersRepo  repository.UsersRepository
	hasher     hash.Hasher
	tokensServ auth.TokensService
}

func NewUsersService(
	usersRepo repository.UsersRepository,
	hasher hash.Hasher,
	tokensServ auth.TokensService,
) *usersService {
	return &usersService{
		usersRepo:  usersRepo,
		hasher:     hasher,
		tokensServ: tokensServ,
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

func (s *usersService) SignIn(ctx context.Context, input UserSignInInput) (auth.Tokens, error) {
	user, err := s.usersRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return auth.Tokens{}, entity.ErrIncorrectEmailOrPassword
		}

		return auth.Tokens{}, errors.Wrapf(err, "could not get user by email %q", input.Email)
	}

	if ok := s.hasher.CheckPasswordHash(input.Password, user.PasswordHash); !ok {
		return auth.Tokens{}, entity.ErrIncorrectEmailOrPassword
	}

	tokens, err := s.tokensServ.CreateTokens(user.ID.Hex())
	if err != nil {
		return auth.Tokens{}, errors.Wrapf(err, "could not create tokens for user with id %s", user.ID.Hex())
	}

	return tokens, nil
}

func (s *usersService) RefreshTokens(_ context.Context, refreshToken string) (auth.Tokens, error) {
	userID, err := s.tokensServ.ParseRefreshToken(refreshToken)
	if err != nil {
		return auth.Tokens{}, errors.Wrapf(err, "could not parse refresh token %q", refreshToken)
	}

	tokens, err := s.tokensServ.CreateTokens(userID)
	if err != nil {
		return auth.Tokens{}, errors.Wrapf(err, "could not create tokens for user with id %s", userID)
	}

	return tokens, nil
}
