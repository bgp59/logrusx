// Collectable logger, (*testing.T).Log style.

// If the test is not running in verbose mode, collect the app logger's output
// and display it JIT at Fatal[f] invocation:

// Typical use:
// func TestSomething(T *testing.T, args... any) {
// 	logger := logrusx_testutils.NewTestCollectableLogger(t, realLogger, nil)
// 	defer logger.RestoreLog()
// }

package logrsux_testutils

import (
	"io"
	"testing"
)

// The interface expected from a collectable logger:
type CollectableLogger interface {
	GetLevel() any
	SetLevel(level any)
	GetOutput() io.Writer
	SetOutput(out io.Writer)
}

type TestCollectableLogger struct {
	logger     CollectableLogger
	savedOut   io.Writer
	savedLevel any
	t          *testing.T
}

func NewTestCollectableLogger(t *testing.T, logger any, level any) *TestCollectableLogger {
	tcl := &TestCollectableLogger{
		t: t,
	}
	if logger, ok := logger.(CollectableLogger); ok && logger != nil {
		if !testing.Verbose() {
			tcl.logger = logger
			tcl.savedOut = logger.GetOutput()
			logger.SetOutput(tcl)
		}
		if level != nil {
			tcl.savedLevel = logger.GetLevel()
			logger.SetLevel(level)
		}
	}
	return tcl
}

func (tcl *TestCollectableLogger) Write(buf []byte) (int, error) {
	n := len(buf)
	if n > 0 && buf[n-1] == '\n' {
		buf = buf[:n-1]
	}
	tcl.t.Log(string(buf))
	return n, nil
}

func (tcl *TestCollectableLogger) RestoreLog() {
	if tcl.logger != nil {
		if tcl.savedOut != nil {
			tcl.logger.SetOutput(tcl.savedOut)
		}
		if tcl.savedLevel != nil {
			tcl.logger.SetLevel(tcl.savedLevel)
		}
	}
}
