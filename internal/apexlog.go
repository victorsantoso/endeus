package internal

import "github.com/apex/log"

var (
	DebugMode = false
)

func SetDebug(debug bool) {
	DebugMode = debug
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}