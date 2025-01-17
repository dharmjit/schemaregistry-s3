// Code generated by MockGen. DO NOT EDIT.
// Source: validator.go
//
// Generated by this command:
//
//	mockgen -package mockschemavalidator -typed=true -source=validator.go -destination ./mocks/validator.go
//

// Package mockschemavalidator is a generated GoMock package.
package mockschemavalidator

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockValidator is a mock of Validator interface.
type MockValidator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorMockRecorder
	isgomock struct{}
}

// MockValidatorMockRecorder is the mock recorder for MockValidator.
type MockValidatorMockRecorder struct {
	mock *MockValidator
}

// NewMockValidator creates a new mock instance.
func NewMockValidator(ctrl *gomock.Controller) *MockValidator {
	mock := &MockValidator{ctrl: ctrl}
	mock.recorder = &MockValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidator) EXPECT() *MockValidatorMockRecorder {
	return m.recorder
}

// ValidateSchema mocks base method.
func (m *MockValidator) ValidateSchema(schema string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSchema", schema)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateSchema indicates an expected call of ValidateSchema.
func (mr *MockValidatorMockRecorder) ValidateSchema(schema any) *MockValidatorValidateSchemaCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSchema", reflect.TypeOf((*MockValidator)(nil).ValidateSchema), schema)
	return &MockValidatorValidateSchemaCall{Call: call}
}

// MockValidatorValidateSchemaCall wrap *gomock.Call
type MockValidatorValidateSchemaCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockValidatorValidateSchemaCall) Return(arg0 error) *MockValidatorValidateSchemaCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockValidatorValidateSchemaCall) Do(f func(string) error) *MockValidatorValidateSchemaCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockValidatorValidateSchemaCall) DoAndReturn(f func(string) error) *MockValidatorValidateSchemaCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
