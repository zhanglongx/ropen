package pkg

import "log"

const (
	// LevelDebug is the debug level
	LevelDebug = iota
	// LevelInfo is the info level
	LevelInfo
)

var level = LevelInfo

// SetLevel sets the log level. XXX: it is not thread safe
func SetLevel(l int) {
	level = l
}

func debug(format string, args ...interface{}) {
	if level <= LevelDebug {
		log.Printf(format, args...)
	}
}
