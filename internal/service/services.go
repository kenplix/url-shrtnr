package service

import (
	"context"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
	"github.com/pkg/errors"
)

type UserSignUpSchema struct {
	Username string
	Email    string
	Password string
}

type UserSignInSchema struct {
	Login    string
	Password string
}

// UsersService is a service for users
//
//go:generate mockery --dir . --name UsersService --output ./mocks
type UsersService interface{}

// AuthService is a service for authorization/authentication
//
//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	SignUp(ctx context.Context, schema UserSignUpSchema) error
	SignIn(ctx context.Context, schema UserSignInSchema) (auth.Tokens, error)
}

type Dependencies struct {
	Repos         *repository.Repositories
	HasherService hash.HasherService
	TokensService auth.TokensService
}

// Services is a collection of all services we have in the project.
type Services struct {
	Users  UsersService
	Auth   AuthService
	Tokens auth.TokensService
}

func NewServices(deps Dependencies) (*Services, error) {
	authServ, err := NewAuthService(deps.Repos.Users, deps.HasherService, deps.TokensService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth service")
	}

	usersServ, err := NewUsersService(deps.Repos.Users)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create users service")
	}

	s := &Services{
		Users:  usersServ,
		Auth:   authServ,
		Tokens: deps.TokensService,
	}

	return s, nil
}
