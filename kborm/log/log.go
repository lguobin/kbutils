package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "[erro] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "[info] ", log.Ldate|log.Ltime|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods
var (
	Print  = infoLog.Println
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// log levels
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel controls log level
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
