package err

import (
	"os"
)

type stringWriter interface {
	WriteString(string) (int, error)
}

var Log = stringWriter(os.Stderr)

func Panic(e error) {
	if e != nil {
		panic(e)
	}
}

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

func Test(e error, t test) {
	if e != nil {
		t.Error(e)
	}
}

type DocumentationError string

func (e DocumentationError) Error() string { return string(e) }

func Issue(issue string) {
	Warn(DocumentationError(issue))
}

func Todo(todo string) {
	Issue("TODO: " + todo)
}

func Depricated(str string) {
	Issue("Depricated: " + str)
}
