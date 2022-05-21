package log

import (
	origlog "log"
	"os"
)

var logger = origlog.New(os.Stderr, "", origlog.LstdFlags|origlog.Llongfile)

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
