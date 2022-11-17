// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/Kenplix/url-shrtnr/pkg/auth"

	mock "github.com/stretchr/testify/mock"

	service "github.com/Kenplix/url-shrtnr/internal/service"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

// SignIn provides a mock function with given fields: ctx, schema
func (_m *AuthService) SignIn(ctx context.Context, schema service.UserSignInSchema) (auth.Tokens, error) {
	ret := _m.Called(ctx, schema)

	var r0 auth.Tokens
	if rf, ok := ret.Get(0).(func(context.Context, service.UserSignInSchema) auth.Tokens); ok {
		r0 = rf(ctx, schema)
	} else {
		r0 = ret.Get(0).(auth.Tokens)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.UserSignInSchema) error); ok {
		r1 = rf(ctx, schema)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUp provides a mock function with given fields: ctx, schema
func (_m *AuthService) SignUp(ctx context.Context, schema service.UserSignUpSchema) error {
	ret := _m.Called(ctx, schema)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, service.UserSignUpSchema) error); ok {
		r0 = rf(ctx, schema)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAuthService interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthService(t mockConstructorTestingTNewAuthService) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
