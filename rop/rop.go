// result.go
package rop

import "fmt"

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](v T) Result[T]         { return Result[T]{value: v} }
func Fail[T any](err error) Result[T] { var zero T; return Result[T]{zero, err} }

func Bind[T any, U any](r Result[T], f func(T) (U, error)) Result[U] {
	if r.err != nil {
		return Fail[U](r.err)
	}
	u, err := f(r.value)
	if err != nil {
		return Fail[U](err)
	}
	return Ok(u)
}

func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.err != nil {
		return Fail[U](r.err)
	}
	return Ok(f(r.value))
}

func Pipe[T any, U any](r Result[T], f func(T) (U, error)) Result[U] {
	return Bind(r, f)
}

func (r Result[T]) Unwrap() (T, error) { return r.value, r.err }
func (r Result[T]) Must() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}
func (r Result[T]) OnError(f func(error)) Result[T] {
	if r.err != nil {
		f(r.err)
	}
	return r
}

func OrElse[T any](r Result[T], f func(error) Result[T]) Result[T] {
	if r.err != nil {
		return f(r.err)
	}
	return r
}

func TeeE[T any](r Result[T], f func(T) error) Result[T] {
	if r.err != nil {
		return r
	}
	if err := f(r.value); err != nil {
		return Fail[T](err)
	}
	return r
}

func Tee[T any](r Result[T], f func(T)) Result[T] {
	if r.err == nil {
		f(r.value)
	}
	return r
}
func (r Result[T]) OnSuccess(f func(T)) Result[T] {
	if r.err == nil {
		f(r.value)
	}
	return r
}

func Try[T any](f func() (T, error)) (res Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			// It's important to assign to the named return variable 'res' here.
			// Otherwise, the Fail result would not be returned.
			res = Fail[T](fmt.Errorf("panic: %v", r))
		}
	}()

	val, err := f()
	if err != nil {
		return Fail[T](err)
	}
	return Ok(val)
}
