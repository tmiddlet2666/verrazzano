// Copyright (C) 2020, 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/verrazzano/verrazzano/application-operator/controllers/loggingscope (interfaces: FluentdManager)

// Package loggingscope is a generated GoMock package.
package loggingscope

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/verrazzano/verrazzano/application-operator/apis/oam/v1alpha1"
)

// MockFluentdManager is a mock of FluentdManager interface
type MockFluentdManager struct {
	ctrl     *gomock.Controller
	recorder *MockFluentdManagerMockRecorder
}

// MockFluentdManagerMockRecorder is the mock recorder for MockFluentdManager
type MockFluentdManagerMockRecorder struct {
	mock *MockFluentdManager
}

// NewMockFluentdManager creates a new mock instance
func NewMockFluentdManager(ctrl *gomock.Controller) *MockFluentdManager {
	mock := &MockFluentdManager{ctrl: ctrl}
	mock.recorder = &MockFluentdManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFluentdManager) EXPECT() *MockFluentdManagerMockRecorder {
	return m.recorder
}

// Apply mocks base method
func (m *MockFluentdManager) Apply(arg0 *LoggingScope, arg1 v1alpha1.QualifiedResourceRelation, arg2 *FluentdPod) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Apply indicates an expected call of Apply
func (mr *MockFluentdManagerMockRecorder) Apply(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockFluentdManager)(nil).Apply), arg0, arg1, arg2)
}

// Remove mocks base method
func (m *MockFluentdManager) Remove(arg0 *LoggingScope, arg1 v1alpha1.QualifiedResourceRelation, arg2 *FluentdPod) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockFluentdManagerMockRecorder) Remove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockFluentdManager)(nil).Remove), arg0, arg1, arg2)
}
