package debug

import (
	"log"
	"os"
)

var debug = os.Getenv("DEBUG") != ""

func Printf(format string, v ...any) {
	if debug {
		log.Printf(format, v...)
	}
}
