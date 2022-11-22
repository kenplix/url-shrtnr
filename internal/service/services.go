package service

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
	"github.com/Kenplix/url-shrtnr/pkg/token"
)

// TokensService provides logic for JWT & Refresh tokens generation, parsing and validation.
//
//go:generate mockery --dir . --name TokensService --output ./mocks
type TokensService interface {
	CreateTokens(ctx context.Context, userID string) (entity.Tokens, error)
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
	Cache          *redis.Client
	Repos          *repository.Repositories
	HasherService  hash.HasherService
	AccessService  token.JWTService
	RefreshService token.JWTService
}

// Services is a collection of all services we have in the project.
type Services struct {
	Tokens TokensService
	Auth   AuthService
	Users  UsersService
}

func NewServices(deps Dependencies) (*Services, error) {
	tokensServ, err := NewTokensService(deps.Cache, deps.AccessService, deps.RefreshService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tokens service")
	}

	authServ, err := NewAuthService(deps.Cache, deps.Repos.Users, deps.HasherService, tokensServ)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth service")
	}

	usersServ, err := NewUsersService(deps.Repos.Users)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create users service")
	}

	s := &Services{
		Tokens: tokensServ,
		Auth:   authServ,
		Users:  usersServ,
	}

	return s, nil
}
