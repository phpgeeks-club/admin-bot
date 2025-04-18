// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotProviderMock is an autogenerated mock type for the BotProvider type
type BotProviderMock struct {
	mock.Mock
}

type BotProviderMock_Expecter struct {
	mock *mock.Mock
}

func (_m *BotProviderMock) EXPECT() *BotProviderMock_Expecter {
	return &BotProviderMock_Expecter{mock: &_m.Mock}
}

// GetChatAdministrators provides a mock function with given fields: chatConfig
func (_m *BotProviderMock) GetChatAdministrators(chatConfig tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error) {
	ret := _m.Called(chatConfig)

	var r0 []tgbotapi.ChatMember
	var r1 error
	if rf, ok := ret.Get(0).(func(tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error)); ok {
		return rf(chatConfig)
	}
	if rf, ok := ret.Get(0).(func(tgbotapi.ChatConfig) []tgbotapi.ChatMember); ok {
		r0 = rf(chatConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]tgbotapi.ChatMember)
		}
	}

	if rf, ok := ret.Get(1).(func(tgbotapi.ChatConfig) error); ok {
		r1 = rf(chatConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BotProviderMock_GetChatAdministrators_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChatAdministrators'
type BotProviderMock_GetChatAdministrators_Call struct {
	*mock.Call
}

// GetChatAdministrators is a helper method to define mock.On call
//   - chatConfig tgbotapi.ChatConfig
func (_e *BotProviderMock_Expecter) GetChatAdministrators(chatConfig interface{}) *BotProviderMock_GetChatAdministrators_Call {
	return &BotProviderMock_GetChatAdministrators_Call{Call: _e.mock.On("GetChatAdministrators", chatConfig)}
}

func (_c *BotProviderMock_GetChatAdministrators_Call) Run(run func(chatConfig tgbotapi.ChatConfig)) *BotProviderMock_GetChatAdministrators_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tgbotapi.ChatConfig))
	})
	return _c
}

func (_c *BotProviderMock_GetChatAdministrators_Call) Return(_a0 []tgbotapi.ChatMember, _a1 error) *BotProviderMock_GetChatAdministrators_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BotProviderMock_GetChatAdministrators_Call) RunAndReturn(run func(tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error)) *BotProviderMock_GetChatAdministrators_Call {
	_c.Call.Return(run)
	return _c
}

// NewMessage provides a mock function with given fields: chatID, text
func (_m *BotProviderMock) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	ret := _m.Called(chatID, text)

	var r0 tgbotapi.MessageConfig
	if rf, ok := ret.Get(0).(func(int64, string) tgbotapi.MessageConfig); ok {
		r0 = rf(chatID, text)
	} else {
		r0 = ret.Get(0).(tgbotapi.MessageConfig)
	}

	return r0
}

// BotProviderMock_NewMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewMessage'
type BotProviderMock_NewMessage_Call struct {
	*mock.Call
}

// NewMessage is a helper method to define mock.On call
//   - chatID int64
//   - text string
func (_e *BotProviderMock_Expecter) NewMessage(chatID interface{}, text interface{}) *BotProviderMock_NewMessage_Call {
	return &BotProviderMock_NewMessage_Call{Call: _e.mock.On("NewMessage", chatID, text)}
}

func (_c *BotProviderMock_NewMessage_Call) Run(run func(chatID int64, text string)) *BotProviderMock_NewMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(string))
	})
	return _c
}

func (_c *BotProviderMock_NewMessage_Call) Return(_a0 tgbotapi.MessageConfig) *BotProviderMock_NewMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BotProviderMock_NewMessage_Call) RunAndReturn(run func(int64, string) tgbotapi.MessageConfig) *BotProviderMock_NewMessage_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: c
func (_m *BotProviderMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	ret := _m.Called(c)

	var r0 tgbotapi.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(tgbotapi.Chattable) (tgbotapi.Message, error)); ok {
		return rf(c)
	}
	if rf, ok := ret.Get(0).(func(tgbotapi.Chattable) tgbotapi.Message); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(tgbotapi.Message)
	}

	if rf, ok := ret.Get(1).(func(tgbotapi.Chattable) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BotProviderMock_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type BotProviderMock_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - c tgbotapi.Chattable
func (_e *BotProviderMock_Expecter) Send(c interface{}) *BotProviderMock_Send_Call {
	return &BotProviderMock_Send_Call{Call: _e.mock.On("Send", c)}
}

func (_c *BotProviderMock_Send_Call) Run(run func(c tgbotapi.Chattable)) *BotProviderMock_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tgbotapi.Chattable))
	})
	return _c
}

func (_c *BotProviderMock_Send_Call) Return(_a0 tgbotapi.Message, _a1 error) *BotProviderMock_Send_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BotProviderMock_Send_Call) RunAndReturn(run func(tgbotapi.Chattable) (tgbotapi.Message, error)) *BotProviderMock_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewBotProviderMock creates a new instance of BotProviderMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBotProviderMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *BotProviderMock {
	mock := &BotProviderMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
