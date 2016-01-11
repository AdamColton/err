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

// Out is the location Warn will write to. The stringWriter interface implements
// WriteString(string) (int, error)
// and is satisfied by *os.File and bufio.Writer
var Out = stringWriter(os.Stderr)
var PanicOnWarn = false

// Panic takes a possible error e and if it is not nil panics
func Panic(e error) {
	if e != nil {
		panic(e)
	}
}

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
			toOut(e)
		}
	}
	return ok
}

// Log takes a possible error e and returns a bool indicating e != nil.
// If e is nil, Log will return true
// If Out is not nil and e is not nil, Log will write e to Out and return false
// If Out is nil and e is not nil, Log will return false
func Log(e error) bool {
	ok := (e == nil)
	if !ok {
		toOut(e)
	}
	return ok
}

// Check takes a possible error e and returns a bool indicating e != nil.
// Same behavior as Log, but with no logging
func Check(e error) bool { return e == nil }

// toOut takes an error and attempts to write it to Out
// if Out is set to nil, false is returned indicating failure
func toOut(e error) {
	if Out == nil {
		return
	}
	errStr := getCallerAndLine(2)
	errStr += e.Error()
	if errStr[len(errStr)-1:] != "\n" {
		errStr += "\n"
	}
	if _, writeError := Out.WriteString(errStr); writeError != nil {
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

// Documentation errors are used to help documentation the development process
// and Warn that features are incomplete, depreciated or any other information
// relevent to active development of current features.
type DocumentationError string

// Error method implements the error interface
func (e DocumentationError) Error() string { return string(e) }

// Issue takes any number of arguments (though generally, one string) and
// concats them to a DocumentationError which is passed to Warn.
// This is useful in documenting sections of code under current development that
// are not ready for release.
func Issue(issue ...interface{}) {
	Warn(DocumentationError(fmt.Sprint(issue...)))
}

// Todo is a wrapper around Issue to specifically call out incomplete items
func Todo(todo ...interface{}) {
	Issue(append([]interface{}{"TODO: "}, todo...))
}

// Depricated is a wrapper around Issue to specifically call out depricated
// code
func Depricated(depricated ...interface{}) {
	Issue(append([]interface{}{"Depricated: "}, depricated...))
}

var DebugEnabled = false
var DebugShowFile = true

// Debug is for printing debug information. It will automatically prepend the
// file name and line. It will print to stdOut, not err.Out By default it is
// disabled, setting DebugEnabled to true will enable it.
func Debug(p ...interface{}) {
	if !DebugEnabled {
		return
	}
	if DebugShowFile {
		fmt.Print(getCallerAndLine(1))
	}
	fmt.Println(p...)
}

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
