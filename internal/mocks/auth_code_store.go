// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/fdelbos/commons/auth"

	mock "github.com/stretchr/testify/mock"
)

// AuthCodeStore is an autogenerated mock type for the CodeStore type
type AuthCodeStore struct {
	mock.Mock
}

// GetCode provides a mock function with given fields: ctx, codeDigest
func (_m *AuthCodeStore) GetCode(ctx context.Context, codeDigest []byte) (*auth.Code, error) {
	ret := _m.Called(ctx, codeDigest)

	var r0 *auth.Code
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) (*auth.Code, error)); ok {
		return rf(ctx, codeDigest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte) *auth.Code); ok {
		r0 = rf(ctx, codeDigest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Code)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, codeDigest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCode provides a mock function with given fields: ctx, code
func (_m *AuthCodeStore) NewCode(ctx context.Context, code *auth.Code) error {
	ret := _m.Called(ctx, code)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *auth.Code) error); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Use provides a mock function with given fields: ctx, codeDigest
func (_m *AuthCodeStore) Use(ctx context.Context, codeDigest []byte) error {
	ret := _m.Called(ctx, codeDigest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) error); ok {
		r0 = rf(ctx, codeDigest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAuthCodeStore creates a new instance of AuthCodeStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthCodeStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthCodeStore {
	mock := &AuthCodeStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
