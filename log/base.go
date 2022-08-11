package log

import (
	origlog "log"
	"os"
)

type LoggerIntf interface {
	Printf(format string, v ...interface{})
}

var logger LoggerIntf = origlog.New(os.Stderr, "", origlog.LstdFlags|origlog.Llongfile)

func SetLogger(newLogger LoggerIntf) func() {
	var backup LoggerIntf
	logger, backup = newLogger, logger
	return func() { logger = backup }
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func TraceMethod(name string) func() {
	Printf("%s START\n", name)
	return func() { Printf("%s END\n", name) }
}
