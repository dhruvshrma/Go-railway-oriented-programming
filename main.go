package main

import (
	"errors"
	"fmt"
	"rop-go/rop"
)

func step1(x int) (string, error) {
	if x < 0 {
		return "", errors.New("negative!")
	}
	y := x * x
	return fmt.Sprintf("v=%d", y), nil
}

func step2(s string) (float64, error) {
	if len(s) == 0 {
		return 0, errors.New("empty")
	}
	return float64(len(s)), nil
}

func bindPattern() {
	fmt.Println("Results using the Bind Pattern")
	res := rop.Ok(42)
	newRes := rop.Bind(res, step1)
	nextRes := rop.Bind(newRes, step2)
	newRes.OnSuccess(func(val string) {
		fmt.Println("step1 success:", val)
	})
	nextRes.OnError(func(err error) {
		fmt.Errorf("step 2 error", err)
	})
}

func pipePattern(input int) {
	fmt.Println("Using pipe Pattern")
	rop.Pipe(
		rop.Pipe(rop.Ok(input), step1),
		step2,
	).OnSuccess(func(val float64) {
		fmt.Println("Final result", val)
	}).OnError(func(err error) {
		fmt.Errorf("Error:", err)
	})
}

func negativeResult(input int) {
	fmt.Println("Negative Result")
	rop.Pipe(rop.Ok(input), step1).
		OnError(func(err error) {
			fmt.Println("Error: Input is negative")
		})
}

func useMapPattern() {
	fmt.Println("Complete Map/Bind example")
	rop.MapExample()
}
func main() {

	// Now by using the pipe pattern
	bindPattern()
	pipePattern(42)
	negativeResult(-10)
	useMapPattern()

	useTeeExample()
	useTeeEExample()
	useOrElseExample()
	useTryExample()
}

func useTeeExample() {
	fmt.Println("\n--- Tee Example ---")
	// Tee with Ok result
	okRes := rop.Ok(100)
	v, err := okRes.Unwrap()
	fmt.Printf("Before Tee (Ok): value=%v, error=%v\n", v, err)
	teeRes1 := rop.Tee(okRes, func(val int) {
		fmt.Printf("Logging value from Tee: %d\n", val)
	})
	v, err = teeRes1.Unwrap()
	fmt.Printf("After Tee (Ok): value=%v, error=%v\n", v, err)

	// Tee with Fail result
	failRes := rop.Fail[int](errors.New("original error for Tee"))
	v, err = failRes.Unwrap()
	fmt.Printf("Before Tee (Fail): value=%v, error=%v\n", v, err)
	teeRes2 := rop.Tee(failRes, func(val int) {
		// This should not be called
		fmt.Printf("Logging value from Tee (should not see this): %d\n", val)
	})
	v, err = teeRes2.Unwrap()
	fmt.Printf("After Tee (Fail): value=%v, error=%v\n", v, err)
}

func useTeeEExample() {
	fmt.Println("\n--- TeeE Example ---")
	// TeeE with Ok result, inner function succeeds
	okRes1 := rop.Ok(200)
	v, err := okRes1.Unwrap()
	fmt.Printf("Before TeeE (Ok, f_ok): value=%v, error=%v\n", v, err)
	teeERes1 := rop.TeeE(okRes1, func(val int) error {
		fmt.Printf("TeeE: Processing value %d, will succeed.\n", val)
		return nil
	})
	v, err = teeERes1.Unwrap()
	fmt.Printf("After TeeE (Ok, f_ok): value=%v, error=%v\n", v, err)

	// TeeE with Ok result, inner function fails
	okRes2 := rop.Ok(300)
	v, err = okRes2.Unwrap()
	fmt.Printf("Before TeeE (Ok, f_fail): value=%v, error=%v\n", v, err)
	teeERes2 := rop.TeeE(okRes2, func(val int) error {
		fmt.Printf("TeeE: Processing value %d, will fail.\n", val)
		return errors.New("error from TeeE's function")
	})
	v, err = teeERes2.Unwrap()
	fmt.Printf("After TeeE (Ok, f_fail): value=%v, error=%v\n", v, err)

	// TeeE with Fail result
	failRes := rop.Fail[int](errors.New("original error for TeeE"))
	v, err = failRes.Unwrap()
	fmt.Printf("Before TeeE (Fail): value=%v, error=%v\n", v, err)
	teeERes3 := rop.TeeE(failRes, func(val int) error {
		// This should not be called
		fmt.Printf("TeeE: Processing value (should not see this): %d\n", val)
		return nil
	})
	v, err = teeERes3.Unwrap()
	fmt.Printf("After TeeE (Fail): value=%v, error=%v\n", v, err)
}

func useOrElseExample() {
	fmt.Println("\n--- OrElse Example ---")
	// OrElse with Fail, recovery returns Ok
	failRes1 := rop.Fail[string](errors.New("first error for OrElse"))
	sVal, err := failRes1.Unwrap()
	fmt.Printf("Before OrElse (Fail, recovery Ok): value=%v, error=%v\n", sVal, err)
	orElseRes1 := rop.OrElse(failRes1, func(e error) rop.Result[string] {
		fmt.Printf("OrElse: Recovering from error '%v', returning Ok with default.\n", e)
		return rop.Ok("default value")
	})
	sVal, err = orElseRes1.Unwrap()
	fmt.Printf("After OrElse (Fail, recovery Ok): value=%v, error=%v\n", sVal, err)

	// OrElse with Fail, recovery returns Fail
	failRes2 := rop.Fail[string](errors.New("second error for OrElse"))
	sVal, err = failRes2.Unwrap()
	fmt.Printf("Before OrElse (Fail, recovery Fail): value=%v, error=%v\n", sVal, err)
	orElseRes2 := rop.OrElse(failRes2, func(e error) rop.Result[string] {
		fmt.Printf("OrElse: Recovering from error '%v', returning another Fail.\n", e)
		return rop.Fail[string](errors.New("error from OrElse's recovery function"))
	})
	sVal, err = orElseRes2.Unwrap()
	fmt.Printf("After OrElse (Fail, recovery Fail): value=%v, error=%v\n", sVal, err)

	// OrElse with Ok
	okRes := rop.Ok("original ok value")
	sVal, err = okRes.Unwrap()
	fmt.Printf("Before OrElse (Ok): value=%v, error=%v\n", sVal, err)
	orElseRes3 := rop.OrElse(okRes, func(e error) rop.Result[string] {
		// This should not be called
		fmt.Printf("OrElse: Recovering (should not see this) from error '%v'.\n", e)
		return rop.Ok("another value")
	})
	sVal, err = orElseRes3.Unwrap()
	fmt.Printf("After OrElse (Ok): value=%v, error=%v\n", sVal, err)
}

func useTryExample() {
	fmt.Println("\n--- Try Example ---")
	// Try with a function that succeeds
	fmt.Println("Trying function that succeeds:")
	tryRes1 := rop.Try(func() (int, error) {
		fmt.Println("Try: Executing function that will succeed.")
		return 42, nil
	})
	v, err := tryRes1.Unwrap()
	fmt.Printf("Result from Try (Success): value=%v, error=%v\n", v, err)

	// Try with a function that returns an error
	fmt.Println("\nTrying function that returns an error:")
	tryRes2 := rop.Try(func() (int, error) {
		fmt.Println("Try: Executing function that will return an error.")
		return 0, errors.New("error from function wrapped by Try")
	})
	v, err = tryRes2.Unwrap()
	fmt.Printf("Result from Try (Returns Error): value=%v, error=%v\n", v, err)

	// Try with a function that panics
	fmt.Println("\nTrying function that panics:")
	tryRes3 := rop.Try(func() (int, error) {
		fmt.Println("Try: Executing function that will panic.")
		panic("oh no, a panic occurred!")
	})
	v, err = tryRes3.Unwrap()
	fmt.Printf("Result from Try (Panics): value=%v, error=%v\n", v, err)
}
