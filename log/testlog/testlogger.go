package testlog

import (
	"testing"

	"github.com/akm/tparser/log"
)

type TestLogger struct {
	t *testing.T
}

var _ log.LoggerIntf = (*TestLogger)(nil)

func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{t: t}
}

func (x *TestLogger) Printf(format string, v ...interface{}) {
	x.t.Logf(format, v...)
}

func Setup(t *testing.T) func() {
	logger := NewTestLogger(t)
	return log.SetLogger(logger)
}
