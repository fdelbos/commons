// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// AuthCodeStore is an autogenerated mock type for the CodeStore type
type AuthCodeStore struct {
	mock.Mock
}

// GetCode provides a mock function with given fields: ctx, code
func (_m *AuthCodeStore) GetCode(ctx context.Context, code string) (string, time.Time, error) {
	ret := _m.Called(ctx, code)

	var r0 string
	var r1 time.Time
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, time.Time, error)); ok {
		return rf(ctx, code)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) time.Time); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, code)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewCode provides a mock function with given fields: ctx, email, code, until
func (_m *AuthCodeStore) NewCode(ctx context.Context, email string, code string, until time.Time) error {
	ret := _m.Called(ctx, email, code, until)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) error); ok {
		r0 = rf(ctx, email, code, until)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Use provides a mock function with given fields: ctx, code
func (_m *AuthCodeStore) Use(ctx context.Context, code string) error {
	ret := _m.Called(ctx, code)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAuthCodeStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthCodeStore creates a new instance of AuthCodeStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthCodeStore(t mockConstructorTestingTNewAuthCodeStore) *AuthCodeStore {
	mock := &AuthCodeStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
