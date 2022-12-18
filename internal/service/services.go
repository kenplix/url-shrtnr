package service

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/repository"
	"github.com/kenplix/url-shrtnr/pkg/hash"
	"github.com/kenplix/url-shrtnr/pkg/token"
)

// JWTService provides logic for JWT & Refresh tokens generation, parsing and validation.
//
//go:generate mockery --dir . --name JWTService --output ./mocks
type JWTService interface {
	CreateTokens(ctx context.Context, userID string) (entity.Tokens, error)
	ProlongTokens(ctx context.Context, userID string)
	ParseAccessToken(token string) (*token.JWTCustomClaims, error)
	ParseRefreshToken(token string) (*token.JWTCustomClaims, error)
	ValidateAccessToken(ctx context.Context, claims *token.JWTCustomClaims) error
	ValidateRefreshToken(ctx context.Context, claims *token.JWTCustomClaims) error
}

type UserSignUpSchema struct {
	Username string
	Email    string
	Password string
}

type UserSignInSchema struct {
	Login    string
	Password string
}

// AuthService is a service for authorization/authentication
//
//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	SignUp(ctx context.Context, schema UserSignUpSchema) error
	SignIn(ctx context.Context, schema UserSignInSchema) (entity.Tokens, error)
	SignOut(ctx context.Context, userID primitive.ObjectID) error
}

// UsersService is a service for users
//
//go:generate mockery --dir . --name UsersService --output ./mocks
type UsersService interface {
	GetByID(ctx context.Context, userID primitive.ObjectID) (entity.User, error)
}

type Dependencies struct {
	Cache            *redis.Client
	Repos            *repository.Repositories
	HasherService    hash.HasherService
	JWTServiceConfig JWTServiceConfig
}

// Services is a collection of all services we have in the project.
type Services struct {
	JWT   JWTService
	Auth  AuthService
	Users UsersService
}

func NewServices(deps Dependencies) (*Services, error) {
	if deps.Repos == nil {
		return nil, errors.New("repositories not provided")
	}

	jwtServ, err := NewJWTService(deps.JWTServiceConfig, deps.Cache)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jwt service")
	}

	authServ, err := NewAuthService(deps.Cache, deps.Repos.Users, deps.HasherService, jwtServ)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth service")
	}

	usersServ, err := NewUsersService(deps.Repos.Users)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create users service")
	}

	s := &Services{
		JWT:   jwtServ,
		Auth:  authServ,
		Users: usersServ,
	}

	return s, nil
}
