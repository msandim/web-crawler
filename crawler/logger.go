package crawler

import (
	"fmt"
	"os"
)

type logger interface {
	logPage(parentURL string, childrenURLs []string)
	logError(msg string)
}

type printer struct{}

func (log *printer) logPage(parentURL string, childrenURLs []string) {
	fmt.Println(". " + parentURL)
	for _, childURL := range childrenURLs {
		fmt.Println("  -> " + childURL)
	}
}

func (log *printer) logError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}
