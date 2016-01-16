package err

import (
	"fmt"
)

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

// Todo is a wrapper around Issue to specifically call ErrOut incomplete items
func Todo(todo ...interface{}) {
	Issue(append([]interface{}{"TODO: "}, todo...))
}

// Deprecated is a wrapper around Issue to specifically call ErrOut depricated
// code
func Deprecated(depricated ...interface{}) {
	Issue(append([]interface{}{"Deprecated: "}, depricated...))
}
