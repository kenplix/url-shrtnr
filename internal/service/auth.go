package service

import (
	"context"
	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
	"github.com/pkg/errors"
	"time"
)

type authService struct {
	usersRepo  repository.UsersRepository
	hasherServ hash.HasherService
	tokensServ auth.TokensService
}

func NewAuthService(
	usersRepo repository.UsersRepository,
	hasherServ hash.HasherService,
	tokensServ auth.TokensService,
) (AuthService, error) {
	if usersRepo == nil {
		return nil, errors.New("users repository not provided")
	}

	if hasherServ == nil {
		return nil, errors.New("hasher service not provided")
	}

	if tokensServ == nil {
		return nil, errors.New("tokens service not provided")
	}

	s := &authService{
		usersRepo:  usersRepo,
		hasherServ: hasherServ,
		tokensServ: tokensServ,
	}

	return s, nil
}

func (s *authService) SignUp(ctx context.Context, schema UserSignUpSchema) error {
	user, err := s.usersRepo.FindByEmail(ctx, schema.Email)
	if err == nil {
		return &entity.ValidationError{
			CoreError: entity.CoreError{
				Code:    errorcode.AlreadyExists,
				Message: "email address already in use by another user",
			},
			Field: "email",
		}
	} else if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		return errors.Wrapf(err, "failed to find user by %q email", schema.Email)
	}

	user, err = s.usersRepo.FindByUsername(ctx, schema.Username)
	if err == nil {
		return &entity.ValidationError{
			CoreError: entity.CoreError{
				Code:    errorcode.AlreadyExists,
				Message: "username is already taken",
			},
			Field: "username",
		}
	} else if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		return errors.Wrapf(err, "failed to find user by %q username", schema.Username)
	}

	passwordHash, err := s.hasherServ.HashPassword(schema.Password)
	if err != nil {
		return errors.Wrapf(err, "failed to hash %q password", schema.Password)
	}

	createdAt := time.Now()

	user = entity.User{
		Username:     schema.Username,
		Email:        schema.Email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}

	err = s.usersRepo.Create(ctx, user)
	if err != nil {
		return errors.Wrapf(err, "failed to create %#v user", user)
	}

	return nil
}

func (s *authService) SignIn(ctx context.Context, schema UserSignInSchema) (auth.Tokens, error) {
	user, err := s.usersRepo.FindByLogin(ctx, schema.Login)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return auth.Tokens{}, entity.ErrIncorrectCredentials
		}

		return auth.Tokens{}, errors.Wrapf(err, "failed to find user by %q login", schema.Login)
	}

	if ok := s.hasherServ.VerifyPassword(schema.Password, user.PasswordHash); !ok {
		return auth.Tokens{}, entity.ErrIncorrectCredentials
	}

	tokens, err := s.tokensServ.CreateTokens(user.ID.Hex())
	if err != nil {
		return auth.Tokens{}, errors.Wrapf(err, "failed to create tokens for user with id %s", user.ID.Hex())
	}

	return tokens, nil
}
