package log

import (
	"log"
	"os"
)

var Enabled bool
var Log func(messages ...any)

func init() {
	logger := log.New(os.Stdout, "", log.Ltime)
	Log = func(messages ...any) {
		if !Enabled {
			return
		}

		logger.Println(messages...)
	}
}
