// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	bytes "bytes"

	mock "github.com/stretchr/testify/mock"
)

// EmailManager is an autogenerated mock type for the EmailManager type
type EmailManager struct {
	mock.Mock
}

// SendEmailBody provides a mock function with given fields: body, subject, tag, recipient
func (_m *EmailManager) SendEmailBody(body bytes.Buffer, subject string, tag string, recipient []string) error {
	ret := _m.Called(body, subject, tag, recipient)

	var r0 error
	if rf, ok := ret.Get(0).(func(bytes.Buffer, string, string, []string) error); ok {
		r0 = rf(body, subject, tag, recipient)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewEmailManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewEmailManager creates a new instance of EmailManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEmailManager(t mockConstructorTestingTNewEmailManager) *EmailManager {
	mock := &EmailManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}