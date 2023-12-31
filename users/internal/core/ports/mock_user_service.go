// Code generated by mockery v2.36.0. DO NOT EDIT.

package ports

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/virsavik/alchemist-template/users/internal/core/domain"
)

// MockUserService is an autogenerated mock type for the UserService type
type MockUserService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *MockUserService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	ret := _m.Called(ctx, user)

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) (domain.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) domain.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, input
func (_m *MockUserService) GetAll(ctx context.Context, input GetUserInput) ([]domain.User, error) {
	ret := _m.Called(ctx, input)

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, GetUserInput) ([]domain.User, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, GetUserInput) []domain.User); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, GetUserInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockUserService creates a new instance of MockUserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserService {
	mock := &MockUserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
