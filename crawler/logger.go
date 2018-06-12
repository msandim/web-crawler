package crawler

import (
	"fmt"
	"os"
)

type logger interface {
	logOutput(msg ...interface{})
	logError(msg ...interface{})
}

type printer struct{}

func (log printer) logOutput(msg ...interface{}) {
	fmt.Println(msg)
}

func (log printer) logError(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg)
}
