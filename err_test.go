package err

import (
	"testing"
)

type testError string

func (e testError) Error() string { return string(e) }

func TestPanic(t *testing.T) {
	defer func() {
		switch recover().(type) {
		case testError:
			return //as expected
		default:
			t.Error("Expected testError")
		}
	}()
	Panic(testError("TestPanic"))
}

type mocker struct {
	buf string
}

func (mock *mocker) WriteString(str string) (int, error) {
	mock.buf = str
	return len(str), nil
}

func (mock *mocker) Error(i ...interface{}) {
	if len(i) != 1 {
		panic("Expected 1, got: ")
	}
	if e, ok := i[0].(error); ok {
		mock.buf = e.Error()
	} else {
		panic("Expected error")
	}
}

func TestWarnNotNil(t *testing.T) {
	mock := &mocker{}
	restore := Log
	Log = mock
	Warn(testError("TestWarnNotNil"))
	if mock.buf != "TestWarnNotNil\n" {
		t.Error("Expected TestWarnNotNil, Got: " + mock.buf)
	}
	Log = restore
}

func TestWarnNil(t *testing.T) {
	restore := Log
	defer func() {
		switch recover().(type) {
		case testError:
			Log = restore
			return //as expected
		default:
			Log = restore
			t.Error("Expected testError")
		}
	}()
	Log = nil
	Warn(testError("TestWarnNil"))
}

func TestTest(t *testing.T) {
	mockTest := &mocker{}
	Test(testError("TestTest"), mockTest)
	if mockTest.buf != "TestTest" {
		t.Error("Exptected TestTest")
	}
}
