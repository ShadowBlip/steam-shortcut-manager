package logger

import (
	"fmt"
	"os"
)

// Print debug messages if debug is enabled
func DebugPrintln(s ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Println(s...)
	}
}
