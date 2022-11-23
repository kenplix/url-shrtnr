package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
)

type authService struct {
	cache      *redis.Client
	usersRepo  repository.UsersRepository
	hasherServ hash.HasherService
	jwtServ    JWTService
}

func NewAuthService(
	cache *redis.Client,
	usersRepo repository.UsersRepository,
	hasherServ hash.HasherService,
	jwtServ JWTService,
) (AuthService, error) {
	if cache == nil {
		return nil, errors.New("cache not provided")
	}

	if usersRepo == nil {
		return nil, errors.New("users repository not provided")
	}

	if hasherServ == nil {
		return nil, errors.New("hasher service not provided")
	}

	if jwtServ == nil {
		return nil, errors.New("jwt service not provided")
	}

	s := &authService{
		cache:      cache,
		usersRepo:  usersRepo,
		hasherServ: hasherServ,
		jwtServ:    jwtServ,
	}

	return s, nil
}

func (s *authService) SignUp(ctx context.Context, schema UserSignUpSchema) error {
	_, err := s.usersRepo.FindByEmail(ctx, schema.Email)
	if err == nil {
		return &entity.ValidationError{
			CoreError: entity.CoreError{
				Code:    errorcode.AlreadyExists,
				Message: "email address already in use by another user",
			},
			Field: "email",
		}
	} else if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		return errors.Wrapf(err, "failed to find user[email:%q]", schema.Email)
	}

	_, err = s.usersRepo.FindByUsername(ctx, schema.Username)
	if err == nil {
		return &entity.ValidationError{
			CoreError: entity.CoreError{
				Code:    errorcode.AlreadyExists,
				Message: "username is already taken",
			},
			Field: "username",
		}
	} else if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		return errors.Wrapf(err, "failed to find user[username:%q]", schema.Username)
	}

	passwordHash, err := s.hasherServ.HashPassword(schema.Password)
	if err != nil {
		return errors.Wrapf(err, "failed to hash %q password", schema.Password)
	}

	now := time.Now()

	user := entity.User{
		Username:     schema.Username,
		Email:        schema.Email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.usersRepo.Create(ctx, user)
	if err != nil {
		return errors.Wrapf(err, "failed to create %+v user", user)
	}

	return nil
}

func (s *authService) SignIn(ctx context.Context, schema UserSignInSchema) (entity.Tokens, error) {
	user, err := s.usersRepo.FindByLogin(ctx, schema.Login)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return entity.Tokens{}, entity.ErrIncorrectCredentials
		}

		return entity.Tokens{}, errors.Wrapf(err, "failed to find user[login:%q]", schema.Login)
	}

	if ok := s.hasherServ.VerifyPassword(schema.Password, user.PasswordHash); !ok {
		return entity.Tokens{}, entity.ErrIncorrectCredentials
	}

	tokens, err := s.jwtServ.CreateTokens(ctx, user.ID.Hex())
	if err != nil {
		return entity.Tokens{}, errors.Wrapf(err, "user[id:%q]: failed to create tokens", user.ID.Hex())
	}

	return tokens, nil
}

func (s *authService) SignOut(ctx context.Context, userID primitive.ObjectID) error {
	userIDHex := userID.Hex()
	tokenKey := tokenCacheKey(userIDHex)

	val, err := s.cache.Del(ctx, tokenKey).Result()
	if err != nil {
		return errors.Wrapf(err, "cache: failed to delete %q key", tokenKey)
	} else if val == 0 {
		return fmt.Errorf("user[id:%q]: already signed out", userIDHex)
	}

	return nil
}
