// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	token "github.com/Kenplix/url-shrtnr/pkg/token"
)

// JWTService is an autogenerated mock type for the JWTService type
type JWTService struct {
	mock.Mock
}

// CreateToken provides a mock function with given fields: id
func (_m *JWTService) CreateToken(id string) (string, string, error) {
	ret := _m.Called(id)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ParseToken provides a mock function with given fields: jwt
func (_m *JWTService) ParseToken(jwt string) (*token.JWTCustomClaims, error) {
	ret := _m.Called(jwt)

	var r0 *token.JWTCustomClaims
	if rf, ok := ret.Get(0).(func(string) *token.JWTCustomClaims); ok {
		r0 = rf(jwt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*token.JWTCustomClaims)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(jwt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenTTL provides a mock function with given fields:
func (_m *JWTService) TokenTTL() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

type mockConstructorTestingTNewJWTService interface {
	mock.TestingT
	Cleanup(func())
}

// NewJWTService creates a new instance of JWTService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJWTService(t mockConstructorTestingTNewJWTService) *JWTService {
	mock := &JWTService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
