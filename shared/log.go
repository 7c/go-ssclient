package shared

import (
	"fmt"
	"log"
	"os"
)

var LoggerConfig struct {
	Verbose bool
}

var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func logf(f string, v ...interface{}) {
	if LoggerConfig.Verbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

type logHelper struct {
	prefix string
}

func (l *logHelper) Write(p []byte) (n int, err error) {
	if LoggerConfig.Verbose {
		logger.Printf("%s%s\n", l.prefix, p)
		return len(p), nil
	}
	return len(p), nil
}

func newLogHelper(prefix string) *logHelper {
	return &logHelper{prefix}
}
