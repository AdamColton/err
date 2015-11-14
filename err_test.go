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
	restore := Out
	Out = mock
	Warn(testError("TestWarnNotNil"))
	if mock.buf != "TestWarnNotNil\n" {
		t.Error("Expected TestWarnNotNil, Got: " + mock.buf)
	}
	Out = restore
}

func TestWarnNil(t *testing.T) {
	restore := Out
	defer func() {
		Out = restore
		switch recover().(type) {
		case testError:
			return //as expected
		default:
			t.Error("Expected testError")
		}
	}()
	Out = nil
	Warn(testError("TestWarnNil"))
}

func TestLogNotNil(t *testing.T) {
	mock := &mocker{}
	restore := Out
	Out = mock
	Log(testError("TestWarnNotNil"))
	if mock.buf != "TestWarnNotNil\n" {
		t.Error("Expected TestWarnNotNil, Got: " + mock.buf)
	}
	Out = restore
}

func TestLogNil(t *testing.T) {
	restore := Out
	defer func() {
		Out = restore
		switch recover().(type) {
		case testError:
			t.Error("Should not have panicked")
		default:
			return //as expected
		}
	}()
	Out = nil
	Log(testError("TestWarnNil"))
}

func TestTest(t *testing.T) {
	mockTest := &mocker{}
	Test(testError("TestTest"), mockTest)
	if mockTest.buf != "TestTest" {
		t.Error("Expected TestTest")
	}
}
