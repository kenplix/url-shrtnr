package usecase

import (
	"context"

	"github.com/Kenplix/url-shrtnr/internal/repository"
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

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// UsersService is a service for users
//go:generate mockery --dir . --name UsersService --output ./mocks
type UsersService interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
}

// Manager is a collection of all services we have in the project.
type Manager struct {
	Users UsersService
}

func NewManager(repos *repository.Repositories) *Manager {
	return &Manager{
		Users: NewUsersService(repos.Users),
	}
}
