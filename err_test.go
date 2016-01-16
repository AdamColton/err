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

func TestWarnNotPanic(t *testing.T) {
	mock := &mocker{}
	restore := ErrOut
	ErrOut = mock
	Warn(testError("TestWarnNotPanic"))
	if l := len(mock.buf); l < 17 || mock.buf[l-17:] != "TestWarnNotPanic\n" {
		t.Error("Expected string ending with 'TestWarnNotPanic', Got: '" + mock.buf + "'")
	}
	ErrOut = restore
}

func TestWarnPanic(t *testing.T) {
	restoreOut := ErrOut
	restorePanic := PanicOnWarn
	defer func() {
		ErrOut = restoreOut
		PanicOnWarn = restorePanic
		switch recover().(type) {
		case testError:
			return //as expected
		default:
			t.Error("Expected testError")
		}
	}()
	ErrOut = nil
	PanicOnWarn = true
	Warn(testError("TestWarnNil"))
}

func TestLog_OutNotNil(t *testing.T) {
	mock := &mocker{}
	restore := ErrOut
	ErrOut = mock
	Log(testError("TestLog_OutNotNil"))
	if l := len(mock.buf); l < 18 || mock.buf[l-18:] != "TestLog_OutNotNil\n" {
		t.Error("Expected ending with 'TestLog_OutNotNil', Got: '" + mock.buf + "'")
	}
	ErrOut = restore
}

func TestLogNil(t *testing.T) {
	restore := ErrOut
	defer func() {
		ErrOut = restore
		switch recover().(type) {
		case testError:
			t.Error("Should not have panicked")
		default:
			return //as expected
		}
	}()
	ErrOut = nil
	Log(testError("TestWarnNil"))
}

func TestTest(t *testing.T) {
	mockTest := &mocker{}
	Test(testError("TestTest"), mockTest)
	if mockTest.buf != "TestTest" {
		t.Error("Expected TestTest")
	}
}
