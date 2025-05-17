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

func main() {
	res := rop.Ok(42)
	newRes := rop.Bind(res, step1)
	nextRes := rop.Bind(newRes, step2)
	newRes.OnSuccess(func(val string) {
		fmt.Println("step1 success:", val)
	})
	nextRes.OnError(func(err error) {
		fmt.Errorf("step 2 error" , err)
	})

}

