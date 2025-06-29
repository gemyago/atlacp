// Code generated by mockery. DO NOT EDIT.

//go:build !release

package app

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockAtlassianAccountsRepository is an autogenerated mock type for the AtlassianAccountsRepository type
type MockAtlassianAccountsRepository struct {
	mock.Mock
}

type MockAtlassianAccountsRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAtlassianAccountsRepository) EXPECT() *MockAtlassianAccountsRepository_Expecter {
	return &MockAtlassianAccountsRepository_Expecter{mock: &_m.Mock}
}

// GetAccountByName provides a mock function with given fields: ctx, name
func (_m *MockAtlassianAccountsRepository) GetAccountByName(ctx context.Context, name string) (*AtlassianAccount, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetAccountByName")
	}

	var r0 *AtlassianAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*AtlassianAccount, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *AtlassianAccount); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AtlassianAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAtlassianAccountsRepository_GetAccountByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccountByName'
type MockAtlassianAccountsRepository_GetAccountByName_Call struct {
	*mock.Call
}

// GetAccountByName is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *MockAtlassianAccountsRepository_Expecter) GetAccountByName(ctx interface{}, name interface{}) *MockAtlassianAccountsRepository_GetAccountByName_Call {
	return &MockAtlassianAccountsRepository_GetAccountByName_Call{Call: _e.mock.On("GetAccountByName", ctx, name)}
}

func (_c *MockAtlassianAccountsRepository_GetAccountByName_Call) Run(run func(ctx context.Context, name string)) *MockAtlassianAccountsRepository_GetAccountByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAtlassianAccountsRepository_GetAccountByName_Call) Return(_a0 *AtlassianAccount, _a1 error) *MockAtlassianAccountsRepository_GetAccountByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAtlassianAccountsRepository_GetAccountByName_Call) RunAndReturn(run func(context.Context, string) (*AtlassianAccount, error)) *MockAtlassianAccountsRepository_GetAccountByName_Call {
	_c.Call.Return(run)
	return _c
}

// GetDefaultAccount provides a mock function with given fields: ctx
func (_m *MockAtlassianAccountsRepository) GetDefaultAccount(ctx context.Context) (*AtlassianAccount, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultAccount")
	}

	var r0 *AtlassianAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*AtlassianAccount, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *AtlassianAccount); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AtlassianAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAtlassianAccountsRepository_GetDefaultAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDefaultAccount'
type MockAtlassianAccountsRepository_GetDefaultAccount_Call struct {
	*mock.Call
}

// GetDefaultAccount is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAtlassianAccountsRepository_Expecter) GetDefaultAccount(ctx interface{}) *MockAtlassianAccountsRepository_GetDefaultAccount_Call {
	return &MockAtlassianAccountsRepository_GetDefaultAccount_Call{Call: _e.mock.On("GetDefaultAccount", ctx)}
}

func (_c *MockAtlassianAccountsRepository_GetDefaultAccount_Call) Run(run func(ctx context.Context)) *MockAtlassianAccountsRepository_GetDefaultAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAtlassianAccountsRepository_GetDefaultAccount_Call) Return(_a0 *AtlassianAccount, _a1 error) *MockAtlassianAccountsRepository_GetDefaultAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAtlassianAccountsRepository_GetDefaultAccount_Call) RunAndReturn(run func(context.Context) (*AtlassianAccount, error)) *MockAtlassianAccountsRepository_GetDefaultAccount_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAtlassianAccountsRepository creates a new instance of MockAtlassianAccountsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAtlassianAccountsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAtlassianAccountsRepository {
	mock := &MockAtlassianAccountsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
