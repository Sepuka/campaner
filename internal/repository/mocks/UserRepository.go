// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/sepuka/campaner/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: userId
func (_m UserRepository) Get(userId int) (*domain.User, error) {
	ret := _m.Called(userId)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(int) *domain.User); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
