package usecase

import (
	"context"

	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
)

type UserSignUpInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UserSignInInput struct {
	Email    string
	Password string
}

// UsersService is a service for users
//
//go:generate mockery --dir . --name UsersService --output ./mocks
type UsersService interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (auth.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (auth.Tokens, error)
}

type Dependencies struct {
	Repos         *repository.Repositories
	HasherService hash.HasherService
	TokensService auth.TokensService
}

// Manager is a collection of all services we have in the project.
type Manager struct {
	Users UsersService
}

func NewManager(deps Dependencies) *Manager {
	return &Manager{
		Users: NewUsersService(deps.Repos.Users, deps.HasherService, deps.TokensService),
	}
}
