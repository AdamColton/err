/*
err is a helper library for even cleaner error handling in Go
*/
package err

import (
	"fmt"
	"os"
	"runtime"
)

type stringWriter interface {
	WriteString(string) (int, error)
}

var Stderr = os.Stderr
var Stdout = os.Stdout

// ErrOut is the location Warn will write to. The stringWriter interface implements
// WriteString(string) (int, error)
// and is satisfied by *os.File and bufio.Writer
var ErrOut stringWriter = os.Stderr

// Panic takes a possible error e and if it is not nil panics
func Panic(e error) {
	if e != nil {
		panic(e)
	}
}

// PanicOnWarn determines the behavior of err.Warn allowing it toggle between
// behaving like err.Log (during development) and err.Panic (when live).
var PanicOnWarn = false

// Warn takes a possible error e and returns a bool indicating e != nil.
// If e is nil, Warn will return true
// If PanicOnWarn is true and e is not nil, Warn call panic(e).
// If PanicOnWarn is true and e is not nil, Warn will log e
func Warn(e error) bool {
	ok := (e == nil)
	if !ok {
		if PanicOnWarn {
			panic(e)
		} else {
			toErrOut(e)
		}
	}
	return ok
}

// Log takes a possible error e and returns a bool indicating e != nil.
// If e is nil, Log will return true
// If ErrOut is not nil and e is not nil, Log will write e to ErrOut and return false
// If ErrOut is nil and e is not nil, Log will return false
func Log(e error) bool {
	ok := (e == nil)
	if !ok {
		toErrOut(e)
	}
	return ok
}

// Check takes a possible error e and returns a bool indicating e != nil.
// Same behavior as Log, but with no logging
func Check(e error) bool { return e == nil }

// toErrOut takes an error and attempts to write it to ErrOut
// if ErrOut is set to nil, false is returned indicating failure
func toErrOut(e error) {
	if ErrOut == nil {
		return
	}
	errStr := getCallerAndLine(2)
	errStr += e.Error()
	if errStr[len(errStr)-1:] != "\n" {
		errStr += "\n"
	}
	if _, writeError := ErrOut.WriteString(errStr); writeError != nil {
		panic(writeError)
	}
}

type test interface {
	Error(...interface{})
}

// Test is for unit testing, the second argument does not need to implement the
// whole testing.T interface, it only requires the Error(...interface{}) method.
// If e is not nil, Test will call t.Error(e)
func Test(e error, t test) {
	if e != nil {
		t.Error(e)
	}
}

// getCallerAndLine attempts to return the file and line number i positions up
// the call stack. It returns only the file name, not the directory. In the case
// that it cannot fetch the information, an empty string is returned.
func getCallerAndLine(i int) string {
	if _, file, line, ok := runtime.Caller(i + 1); ok {
		for i := len(file) - 1; i >= 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}
		return fmt.Sprint(file, ":", line, ": ")
	}
	return ""
}
