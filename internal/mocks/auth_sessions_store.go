// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	auth "github.com/fdelbos/commons/auth"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// AuthSessionsStore is an autogenerated mock type for the SessionsStore type
type AuthSessionsStore struct {
	mock.Mock
}

// All provides a mock function with given fields: userID
func (_m *AuthSessionsStore) All(userID uuid.UUID) ([]*auth.Session, error) {
	ret := _m.Called(userID)

	var r0 []*auth.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) ([]*auth.Session, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) []*auth.Session); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*auth.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields: sessionID
func (_m *AuthSessionsStore) Close(sessionID string) error {
	ret := _m.Called(sessionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(sessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: sessionID
func (_m *AuthSessionsStore) Get(sessionID string) (*auth.Session, error) {
	ret := _m.Called(sessionID)

	var r0 *auth.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*auth.Session, error)); ok {
		return rf(sessionID)
	}
	if rf, ok := ret.Get(0).(func(string) *auth.Session); ok {
		r0 = rf(sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sessionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// New provides a mock function with given fields: session
func (_m *AuthSessionsStore) New(session auth.Session) error {
	ret := _m.Called(session)

	var r0 error
	if rf, ok := ret.Get(0).(func(auth.Session) error); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAuthSessionsStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthSessionsStore creates a new instance of AuthSessionsStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthSessionsStore(t mockConstructorTestingTNewAuthSessionsStore) *AuthSessionsStore {
	mock := &AuthSessionsStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}