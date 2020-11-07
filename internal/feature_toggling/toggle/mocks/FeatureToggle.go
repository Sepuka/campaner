// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/sepuka/campaner/internal/feature_toggling/domain"
	mock "github.com/stretchr/testify/mock"
)

// FeatureToggle is an autogenerated mock type for the FeatureToggle type
type FeatureToggle struct {
	mock.Mock
}

// IsEnabled provides a mock function with given fields: userId, feature
func (_m FeatureToggle) IsEnabled(userId int, feature domain.FeatureName) bool {
	ret := _m.Called(userId, feature)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int, domain.FeatureName) bool); ok {
		r0 = rf(userId, feature)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}