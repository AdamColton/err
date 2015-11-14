/**
 * err is a helper library for even cleaner error handling in g
 */
package err

import (
	"fmt"
	"os"
)

type stringWriter interface {
	WriteString(string) (int, error)
}

// Log is the location Warn will write to. The stringWriter interface implements
// WriteString(string) (int, error)
// and is satisfied by *os.File and bufio.Writer
var Log = stringWriter(os.Stderr)

// Panic takes a possible error e and if it is not nil panics
func Panic(e error) {
	if e != nil {
		panic(e)
	}
}

// Warn takes a possible error e and returns a bool indicating e != nil.
// If Log is nil and e is not nil, Warn will call panic(e).
// If Log is not and e is not nil, Warn will write e to Log.
func Warn(e error) bool {
	if e != nil {
		if Log != nil {
			errStr := e.Error()
			if errStr[len(errStr)-1:] != "\n" {
				errStr += "\n"
			}
			if _, writeError := Log.WriteString(errStr); writeError != nil {
				panic(writeError)
			}
		} else {
			panic(e)
		}
		return false
	}
	return true
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
