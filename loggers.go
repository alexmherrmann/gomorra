package gomorra

import (
	"log"
	"os"
)

var debugLogger *log.Logger = log.New(os.Stderr, "DEBUG", log.LstdFlags)
func D() *log.Logger {
	return debugLogger
}