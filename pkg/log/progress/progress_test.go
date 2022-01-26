// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package progress

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/verrazzano/verrazzano/pkg/log"
	kzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
	"time"

	"go.uber.org/zap"
)

type fakeLogger struct {
	expectedMsg string
	actualMsg   string
	count       int
}

// TestLog tests the ProgressLogger function periodic logging
// GIVEN a ProgressLogger with a frequency of 3 seconds
// WHEN log is called 5 times in 5 seconds to log the same message
// THEN ensure that 2 messages are logged
func TestLog(t *testing.T) {
	msg := "test1"
	logger := fakeLogger{expectedMsg: msg}
	const rKey = "testns/test"
	rl := EnsureLogContext(rKey, &logger)
	l := rl.EnsureProgressLogger("comp1").SetFrequency(3)

	// 5 calls to log should result in only 2 log messages being written
	// since the frequency is 3 secs
	for i := 0; i < 5; i++ {
		l.Progress(msg)
		time.Sleep(time.Duration(1) * time.Second)
	}
	assert.Equal(t, 2, logger.count)
	assert.Equal(t, logger.actualMsg, logger.expectedMsg)
	DeleteLogContext(rKey)
}

// TestLogNewMsg tests the ProgressLogger function periodic logging
// GIVEN a ProgressLogger with a frequency of 2 seconds
// WHEN log is called 5 times with 2 different message
// THEN ensure that 2 messages are logged
func TestLogNewMsg(t *testing.T) {
	msg := "test1"
	msg2 := "test2"
	logger := fakeLogger{expectedMsg: msg}
	const rKey = "testns/test2"
	rl := EnsureLogContext(rKey, &logger)
	l := rl.EnsureProgressLogger("comp1").SetFrequency(2)

	// Calls to log should result in only 2 log messages being written
	l.Progress(msg)
	l.Progress(msg)
	l.Progress(msg)
	l.Progress(msg2)
	l.Progress(msg2)
	assert.Equal(t, 2, logger.count)
	assert.Equal(t, logger.actualMsg, msg2)
	DeleteLogContext(rKey)
}

// TestLogFormat tests the ProgressLogger function message formatting
// GIVEN a ProgressLogger
// WHEN log.Infof is called with a string and a template
// THEN ensure that the message is formatted correctly and logged
func TestLogFormat(t *testing.T) {
	template := "test %s"
	inStr := "foo"
	logger := fakeLogger{}
	logger.expectedMsg = fmt.Sprintf(template, inStr)
	const rKey = "testns/test3"
	rl := EnsureLogContext(rKey, &logger)
	l := rl.EnsureProgressLogger("comp1")
	l.Progressf(template, inStr)
	assert.Equal(t, 1, logger.count)
	assert.Equal(t, logger.actualMsg, logger.expectedMsg)
	DeleteLogContext(rKey)
}

// TestDefault tests the DefaultProgressLogger
// GIVEN a DefaultProgressLogger
// WHEN log.Infof is called with a string and a template
// THEN ensure that the message is formatted correctly and logged
func TestDefault(t *testing.T) {
	template := "test %s"
	inStr := "foo"
	logger := fakeLogger{}
	logger.expectedMsg = fmt.Sprintf(template, inStr)
	const rKey = "testns/test3"
	l := EnsureLogContext(rKey, &logger).DefaultProgressLogger()
	l.Progressf(template, inStr)
	assert.Equal(t, 1, logger.count)
	assert.Equal(t, logger.actualMsg, logger.expectedMsg)
	DeleteLogContext(rKey)
}

// TestMultipleContexts tests the EnsureLogContext and DeleteLogContext
// WHEN EnsureLogContext is called multiple times
// THEN ensure that the context map has an entry for each context and that
//   the context map is empty when they the contexts are deleted
func TestMultipleContexts(t *testing.T) {
	const rKey1 = "k1"
	const rKey2 = "k2"
	logger1 := fakeLogger{}
	logger2 := fakeLogger{}
	c1 := EnsureLogContext(rKey1, &logger1)
	c2 := EnsureLogContext(rKey2, &logger2)

	assert.Equal(t, 2, len(LogContextMap))
	c1Actual := LogContextMap[rKey1]
	assert.Equal(t, c1, c1Actual)
	c2Actual := LogContextMap[rKey2]
	assert.Equal(t, c2, c2Actual)
	DeleteLogContext(rKey1)
	DeleteLogContext(rKey2)
	assert.Equal(t, 0, len(LogContextMap))
}

// TestZap tests the zap SugaredLogger
// GIVEN a zap SugaredLogger
// WHEN EnsureLogContext is called with the SugaredLogger
// THEN ensure that the ProgressMessage can be called
func TestZap(t *testing.T) {
	testOpts := kzap.Options{}
	testOpts.Development = true
	testOpts.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	log.InitLogs(testOpts)
	const rKey = "testns/test3"
	l := EnsureLogContext(rKey, zap.S()).DefaultProgressLogger()
	l.Progress("testmsg")
	DeleteLogContext(rKey)
}

func (l *fakeLogger) Info(args ...interface{}) {
	s := fmt.Sprint(args...)
	l.actualMsg = s
	l.count = l.count + 1
	fmt.Println(s)
}

// Infof formats a message and logs it
func (l *fakeLogger) Infof(template string, args ...interface{}) {
	s := fmt.Sprintf(template, args...)
	l.Info(s)
}

// Debug is a wrapper for SugaredLogger Debug
func (l *fakeLogger) Debug(args ...interface{}) {
}

// Debugf is a wrapper for SugaredLogger Debugf
func (l *fakeLogger) Debugf(template string, args ...interface{}) {
}

// Error is a wrapper for SugaredLogger Error
func (l *fakeLogger) Error(args ...interface{}) {
}

// Errorf is a wrapper for SugaredLogger Errorf
func (l *fakeLogger) Errorf(template string, args ...interface{}) {
}
