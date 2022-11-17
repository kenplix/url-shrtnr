// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/Kenplix/url-shrtnr/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// UsersRepository is an autogenerated mock type for the UsersRepository type
type UsersRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *UsersRepository) Create(ctx context.Context, user entity.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *UsersRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	ret := _m.Called(ctx, email)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByLogin provides a mock function with given fields: ctx, login
func (_m *UsersRepository) FindByLogin(ctx context.Context, login string) (entity.User, error) {
	ret := _m.Called(ctx, login)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.User); ok {
		r0 = rf(ctx, login)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUsername provides a mock function with given fields: ctx, username
func (_m *UsersRepository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	ret := _m.Called(ctx, username)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.User); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUsersRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUsersRepository creates a new instance of UsersRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUsersRepository(t mockConstructorTestingTNewUsersRepository) *UsersRepository {
	mock := &UsersRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
