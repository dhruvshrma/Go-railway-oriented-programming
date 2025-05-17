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

}
