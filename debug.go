package err

import (
	"fmt"
	"os"
)

// ErrOut is the location Warn will write to. The stringWriter interface implements
// WriteString(string) (int, error)
// and is satisfied by *os.File and bufio.Writer
var DebugOut stringWriter = os.Stdout

// DebugShowFile, when true, will preced any debug output with the file and line
// of the caller.
var DebugShowFile = true

// Debug is for printing debug information. It will automatically prepend the
// file name and line. It will print to stdOut, not err.Out By default it is
// disabled, setting DebugEnabled to true will enable it.
func Debug(p ...interface{}) {
	if DebugOut != nil {
		return
	}

	if DebugShowFile {
		if _, writeError := ErrOut.WriteString(getCallerAndLine(1)); writeError != nil {
			panic(writeError)
		}
	}
	if _, writeError := ErrOut.WriteString(fmt.Sprintln(p...)); writeError != nil {
		panic(writeError)
	}
}

// DebugEnabled is just a shorthand to check if Debugging output is enabled
func DebugEnabled() bool {
	return DebugOut != nil
}
