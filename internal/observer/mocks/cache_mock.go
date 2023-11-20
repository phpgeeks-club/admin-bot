// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CacheMock is an autogenerated mock type for the Cache type
type CacheMock struct {
	mock.Mock
}

type CacheMock_Expecter struct {
	mock *mock.Mock
}

func (_m *CacheMock) EXPECT() *CacheMock_Expecter {
	return &CacheMock_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: key
func (_m *CacheMock) Get(key string) ([]tgbotapi.ChatMember, bool) {
	ret := _m.Called(key)

	var r0 []tgbotapi.ChatMember
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) ([]tgbotapi.ChatMember, bool)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) []tgbotapi.ChatMember); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]tgbotapi.ChatMember)
		}
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// CacheMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type CacheMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *CacheMock_Expecter) Get(key interface{}) *CacheMock_Get_Call {
	return &CacheMock_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *CacheMock_Get_Call) Run(run func(key string)) *CacheMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CacheMock_Get_Call) Return(value []tgbotapi.ChatMember, ok bool) *CacheMock_Get_Call {
	_c.Call.Return(value, ok)
	return _c
}

func (_c *CacheMock_Get_Call) RunAndReturn(run func(string) ([]tgbotapi.ChatMember, bool)) *CacheMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: key, value
func (_m *CacheMock) Set(key string, value []tgbotapi.ChatMember) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []tgbotapi.ChatMember) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheMock_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type CacheMock_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - key string
//   - value []tgbotapi.ChatMember
func (_e *CacheMock_Expecter) Set(key interface{}, value interface{}) *CacheMock_Set_Call {
	return &CacheMock_Set_Call{Call: _e.mock.On("Set", key, value)}
}

func (_c *CacheMock_Set_Call) Run(run func(key string, value []tgbotapi.ChatMember)) *CacheMock_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]tgbotapi.ChatMember))
	})
	return _c
}

func (_c *CacheMock_Set_Call) Return(_a0 error) *CacheMock_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheMock_Set_Call) RunAndReturn(run func(string, []tgbotapi.ChatMember) error) *CacheMock_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewCacheMock creates a new instance of CacheMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCacheMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *CacheMock {
	mock := &CacheMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
