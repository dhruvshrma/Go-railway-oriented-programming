// result.go
package rop

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](v T) Result[T]         { return Result[T]{value: v} }
func Fail[T any](err error) Result[T] { var zero T; return Result[T]{zero, err} }

// func (r Result[T]) Bind[U any](f func(T) (U, error)) Result[U] {
//     if r.err != nil {
//         return Fail[U](r.err)
//     }
//     u, err := f(r.value)
//     if err != nil {
//         return Fail[U](err)
//     }
//     return Ok(u)
// }

// func (r Result[T]) Then[U any](f func(T) (U, error)) func() Result[U] {
//     return func() Result[U] {
//         if r.err != nil {
//             return Fail[U](r.err)
//         }
//         u, err := f(r.value)
//         if err != nil {
//             return Fail[U](err)
//         }
//         return Ok(u)
//     }
// }

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
func (r Result[T]) OnSuccess(f func(T)) Result[T] {
	if r.err == nil {
		f(r.value)
	}
	return r
}
