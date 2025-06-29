// Code generated by mockery. DO NOT EDIT.

//go:build !release

package controllers

import (
	context "context"

	app "github.com/gemyago/atlacp/internal/app"
	bitbucket "github.com/gemyago/atlacp/internal/services/bitbucket"

	mock "github.com/stretchr/testify/mock"
)

// MockbitbucketService is an autogenerated mock type for the bitbucketService type
type MockbitbucketService struct {
	mock.Mock
}

type MockbitbucketService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockbitbucketService) EXPECT() *MockbitbucketService_Expecter {
	return &MockbitbucketService_Expecter{mock: &_m.Mock}
}

// ApprovePR provides a mock function with given fields: ctx, params
func (_m *MockbitbucketService) ApprovePR(ctx context.Context, params app.BitbucketApprovePRParams) (*bitbucket.Participant, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for ApprovePR")
	}

	var r0 *bitbucket.Participant
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketApprovePRParams) (*bitbucket.Participant, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketApprovePRParams) *bitbucket.Participant); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bitbucket.Participant)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, app.BitbucketApprovePRParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockbitbucketService_ApprovePR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApprovePR'
type MockbitbucketService_ApprovePR_Call struct {
	*mock.Call
}

// ApprovePR is a helper method to define mock.On call
//   - ctx context.Context
//   - params app.BitbucketApprovePRParams
func (_e *MockbitbucketService_Expecter) ApprovePR(ctx interface{}, params interface{}) *MockbitbucketService_ApprovePR_Call {
	return &MockbitbucketService_ApprovePR_Call{Call: _e.mock.On("ApprovePR", ctx, params)}
}

func (_c *MockbitbucketService_ApprovePR_Call) Run(run func(ctx context.Context, params app.BitbucketApprovePRParams)) *MockbitbucketService_ApprovePR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.BitbucketApprovePRParams))
	})
	return _c
}

func (_c *MockbitbucketService_ApprovePR_Call) Return(_a0 *bitbucket.Participant, _a1 error) *MockbitbucketService_ApprovePR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockbitbucketService_ApprovePR_Call) RunAndReturn(run func(context.Context, app.BitbucketApprovePRParams) (*bitbucket.Participant, error)) *MockbitbucketService_ApprovePR_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePR provides a mock function with given fields: ctx, params
func (_m *MockbitbucketService) CreatePR(ctx context.Context, params app.BitbucketCreatePRParams) (*bitbucket.PullRequest, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for CreatePR")
	}

	var r0 *bitbucket.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketCreatePRParams) (*bitbucket.PullRequest, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketCreatePRParams) *bitbucket.PullRequest); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bitbucket.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, app.BitbucketCreatePRParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockbitbucketService_CreatePR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePR'
type MockbitbucketService_CreatePR_Call struct {
	*mock.Call
}

// CreatePR is a helper method to define mock.On call
//   - ctx context.Context
//   - params app.BitbucketCreatePRParams
func (_e *MockbitbucketService_Expecter) CreatePR(ctx interface{}, params interface{}) *MockbitbucketService_CreatePR_Call {
	return &MockbitbucketService_CreatePR_Call{Call: _e.mock.On("CreatePR", ctx, params)}
}

func (_c *MockbitbucketService_CreatePR_Call) Run(run func(ctx context.Context, params app.BitbucketCreatePRParams)) *MockbitbucketService_CreatePR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.BitbucketCreatePRParams))
	})
	return _c
}

func (_c *MockbitbucketService_CreatePR_Call) Return(_a0 *bitbucket.PullRequest, _a1 error) *MockbitbucketService_CreatePR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockbitbucketService_CreatePR_Call) RunAndReturn(run func(context.Context, app.BitbucketCreatePRParams) (*bitbucket.PullRequest, error)) *MockbitbucketService_CreatePR_Call {
	_c.Call.Return(run)
	return _c
}

// MergePR provides a mock function with given fields: ctx, params
func (_m *MockbitbucketService) MergePR(ctx context.Context, params app.BitbucketMergePRParams) (*bitbucket.PullRequest, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for MergePR")
	}

	var r0 *bitbucket.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketMergePRParams) (*bitbucket.PullRequest, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketMergePRParams) *bitbucket.PullRequest); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bitbucket.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, app.BitbucketMergePRParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockbitbucketService_MergePR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MergePR'
type MockbitbucketService_MergePR_Call struct {
	*mock.Call
}

// MergePR is a helper method to define mock.On call
//   - ctx context.Context
//   - params app.BitbucketMergePRParams
func (_e *MockbitbucketService_Expecter) MergePR(ctx interface{}, params interface{}) *MockbitbucketService_MergePR_Call {
	return &MockbitbucketService_MergePR_Call{Call: _e.mock.On("MergePR", ctx, params)}
}

func (_c *MockbitbucketService_MergePR_Call) Run(run func(ctx context.Context, params app.BitbucketMergePRParams)) *MockbitbucketService_MergePR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.BitbucketMergePRParams))
	})
	return _c
}

func (_c *MockbitbucketService_MergePR_Call) Return(_a0 *bitbucket.PullRequest, _a1 error) *MockbitbucketService_MergePR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockbitbucketService_MergePR_Call) RunAndReturn(run func(context.Context, app.BitbucketMergePRParams) (*bitbucket.PullRequest, error)) *MockbitbucketService_MergePR_Call {
	_c.Call.Return(run)
	return _c
}

// ReadPR provides a mock function with given fields: ctx, params
func (_m *MockbitbucketService) ReadPR(ctx context.Context, params app.BitbucketReadPRParams) (*bitbucket.PullRequest, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for ReadPR")
	}

	var r0 *bitbucket.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketReadPRParams) (*bitbucket.PullRequest, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketReadPRParams) *bitbucket.PullRequest); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bitbucket.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, app.BitbucketReadPRParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockbitbucketService_ReadPR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadPR'
type MockbitbucketService_ReadPR_Call struct {
	*mock.Call
}

// ReadPR is a helper method to define mock.On call
//   - ctx context.Context
//   - params app.BitbucketReadPRParams
func (_e *MockbitbucketService_Expecter) ReadPR(ctx interface{}, params interface{}) *MockbitbucketService_ReadPR_Call {
	return &MockbitbucketService_ReadPR_Call{Call: _e.mock.On("ReadPR", ctx, params)}
}

func (_c *MockbitbucketService_ReadPR_Call) Run(run func(ctx context.Context, params app.BitbucketReadPRParams)) *MockbitbucketService_ReadPR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.BitbucketReadPRParams))
	})
	return _c
}

func (_c *MockbitbucketService_ReadPR_Call) Return(_a0 *bitbucket.PullRequest, _a1 error) *MockbitbucketService_ReadPR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockbitbucketService_ReadPR_Call) RunAndReturn(run func(context.Context, app.BitbucketReadPRParams) (*bitbucket.PullRequest, error)) *MockbitbucketService_ReadPR_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePR provides a mock function with given fields: ctx, params
func (_m *MockbitbucketService) UpdatePR(ctx context.Context, params app.BitbucketUpdatePRParams) (*bitbucket.PullRequest, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePR")
	}

	var r0 *bitbucket.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketUpdatePRParams) (*bitbucket.PullRequest, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, app.BitbucketUpdatePRParams) *bitbucket.PullRequest); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bitbucket.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, app.BitbucketUpdatePRParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockbitbucketService_UpdatePR_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePR'
type MockbitbucketService_UpdatePR_Call struct {
	*mock.Call
}

// UpdatePR is a helper method to define mock.On call
//   - ctx context.Context
//   - params app.BitbucketUpdatePRParams
func (_e *MockbitbucketService_Expecter) UpdatePR(ctx interface{}, params interface{}) *MockbitbucketService_UpdatePR_Call {
	return &MockbitbucketService_UpdatePR_Call{Call: _e.mock.On("UpdatePR", ctx, params)}
}

func (_c *MockbitbucketService_UpdatePR_Call) Run(run func(ctx context.Context, params app.BitbucketUpdatePRParams)) *MockbitbucketService_UpdatePR_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(app.BitbucketUpdatePRParams))
	})
	return _c
}

func (_c *MockbitbucketService_UpdatePR_Call) Return(_a0 *bitbucket.PullRequest, _a1 error) *MockbitbucketService_UpdatePR_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockbitbucketService_UpdatePR_Call) RunAndReturn(run func(context.Context, app.BitbucketUpdatePRParams) (*bitbucket.PullRequest, error)) *MockbitbucketService_UpdatePR_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockbitbucketService creates a new instance of MockbitbucketService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockbitbucketService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockbitbucketService {
	mock := &MockbitbucketService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
