package logger

import (
	l "log"
	"os"
)

var Logger *l.Logger

func init() {
	Logger = l.New(os.Stdout, "legogen", l.Lshortfile|l.Ltime|l.Ldate|l.Lmicroseconds)
}
