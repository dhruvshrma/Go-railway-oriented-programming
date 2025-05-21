package rop

import (
	"errors"
	"testing"
)

func TestTee_OkResult(t *testing.T) {
	var called bool
	var valuePassed int
	f := func(v int) {
		called = true
		valuePassed = v
	}

	result := Ok(10)
	returnedResult := Tee(result, f)

	if !called {
		t.Error("Expected function f to be called, but it was not")
	}
	if valuePassed != 10 {
		t.Errorf("Expected function f to be called with value 10, but got %d", valuePassed)
	}
	if returnedResult.err != nil {
		t.Errorf("Expected Ok result to be returned, but got Fail: %v", returnedResult.err)
	}
	if returnedResult.value != 10 {
		t.Errorf("Expected Ok result with value 10, but got value %d", returnedResult.value)
	}
}

func TestTee_FailResult(t *testing.T) {
	var called bool
	f := func(v int) {
		called = true
	}

	err := errors.New("test error")
	result := Fail[int](err)
	returnedResult := Tee(result, f)

	if called {
		t.Error("Expected function f not to be called, but it was")
	}
	if returnedResult.err == nil {
		t.Error("Expected Fail result to be returned, but got Ok")
	}
	if returnedResult.err != err {
		t.Errorf("Expected Fail result with error '%v', but got error '%v'", err, returnedResult.err)
	}
}

func TestTry_Success(t *testing.T) {
	f := func() (int, error) {
		return 123, nil
	}

	result := Try(f)

	if result.err != nil {
		t.Errorf("Expected Ok result, but got error: %v", result.err)
	}
	if result.value != 123 {
		t.Errorf("Expected value 123, but got %v", result.value)
	}
}

func TestTry_ReturnsError(t *testing.T) {
	expectedErr := errors.New("test error from f")
	f := func() (int, error) {
		return 0, expectedErr
	}

	result := Try(f)

	if result.err == nil {
		t.Error("Expected Fail result, but got Ok")
	}
	if result.err != expectedErr {
		t.Errorf("Expected error '%v', but got '%v'", expectedErr, result.err)
	}
}

func TestTry_Panics(t *testing.T) {
	f := func() (int, error) {
		panic("test panic")
	}

	result := Try(f)

	if result.err == nil {
		t.Error("Expected Fail result from panic, but got Ok")
	}
	expectedPanicMsg := "panic: test panic"
	if result.err.Error() != expectedPanicMsg {
		t.Errorf("Expected panic message '%s', but got '%s'", expectedPanicMsg, result.err.Error())
	}
}

func TestOrElse_OkResult(t *testing.T) {
	var fCalled bool
	f := func(e error) Result[int] {
		fCalled = true
		return Fail[int](e) // Should not be called
	}

	originalResult := Ok(123)
	returnedResult := OrElse(originalResult, f)

	if fCalled {
		t.Error("Expected function f NOT to be called, but it was")
	}
	if returnedResult.err != nil {
		t.Errorf("Expected Ok result, but got error: %v", returnedResult.err)
	}
	if returnedResult.value != 123 {
		t.Errorf("Expected value 123, but got %v", returnedResult.value)
	}
}

func TestOrElse_FailResult_FReturnsOk(t *testing.T) {
	var fCalled bool
	var errPassed error
	f := func(e error) Result[int] {
		fCalled = true
		errPassed = e
		return Ok(456)
	}

	originalErr := errors.New("original error")
	originalResult := Fail[int](originalErr)
	returnedResult := OrElse(originalResult, f)

	if !fCalled {
		t.Error("Expected function f to be called, but it was not")
	}
	if errPassed != originalErr {
		t.Errorf("Expected f to be called with error '%v', but got '%v'", originalErr, errPassed)
	}
	if returnedResult.err != nil {
		t.Errorf("Expected Ok result from f, but got error: %v", returnedResult.err)
	}
	if returnedResult.value != 456 {
		t.Errorf("Expected value 456 from f, but got %v", returnedResult.value)
	}
}

func TestOrElse_FailResult_FReturnsFail(t *testing.T) {
	var fCalled bool
	var errPassed error
	fErr := errors.New("error from f")
	f := func(e error) Result[int] {
		fCalled = true
		errPassed = e
		return Fail[int](fErr)
	}

	originalErr := errors.New("original error")
	originalResult := Fail[int](originalErr)
	returnedResult := OrElse(originalResult, f)

	if !fCalled {
		t.Error("Expected function f to be called, but it was not")
	}
	if errPassed != originalErr {
		t.Errorf("Expected f to be called with error '%v', but got '%v'", originalErr, errPassed)
	}
	if returnedResult.err == nil {
		t.Error("Expected Fail result from f, but got Ok")
	}
	if returnedResult.err != fErr {
		t.Errorf("Expected error '%v' from f, but got '%v'", fErr, returnedResult.err)
	}
}

func TestTeeE_OkResult_FReturnsNil(t *testing.T) {
	var called bool
	var valuePassed int
	f := func(v int) error {
		called = true
		valuePassed = v
		return nil
	}

	result := Ok(10)
	returnedResult := TeeE(result, f)

	if !called {
		t.Error("Expected function f to be called, but it was not")
	}
	if valuePassed != 10 {
		t.Errorf("Expected function f to be called with value 10, but got %d", valuePassed)
	}
	if returnedResult.err != nil {
		t.Errorf("Expected Ok result to be returned, but got Fail: %v", returnedResult.err)
	}
	if returnedResult.value != 10 {
		t.Errorf("Expected Ok result with value 10, but got value %d", returnedResult.value)
	}
}

func TestTeeE_OkResult_FReturnsError(t *testing.T) {
	var called bool
	var valuePassed int
	errF := errors.New("error from f")
	f := func(v int) error {
		called = true
		valuePassed = v
		return errF
	}

	result := Ok(10)
	returnedResult := TeeE(result, f)

	if !called {
		t.Error("Expected function f to be called, but it was not")
	}
	if valuePassed != 10 {
		t.Errorf("Expected function f to be called with value 10, but got %d", valuePassed)
	}
	if returnedResult.err == nil {
		t.Error("Expected Fail result to be returned, but got Ok")
	}
	if returnedResult.err != errF {
		t.Errorf("Expected Fail result with error '%v', but got error '%v'", errF, returnedResult.err)
	}
}

func TestTeeE_FailResult(t *testing.T) {
	var called bool
	f := func(v int) error {
		called = true
		return nil
	}

	err := errors.New("test error")
	result := Fail[int](err)
	returnedResult := TeeE(result, f)

	if called {
		t.Error("Expected function f not to be called, but it was")
	}
	if returnedResult.err == nil {
		t.Error("Expected Fail result to be returned, but got Ok")
	}
	if returnedResult.err != err {
		t.Errorf("Expected Fail result with error '%v', but got error '%v'", err, returnedResult.err)
	}
}
