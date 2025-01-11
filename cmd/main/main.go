package main

import (
	"bitbucket.org/drbolle/go-boli/internal/interpreter"
	"fmt"
)

func main() {
	code := `(+ 41 1)`
	value, err := interpreter.Run(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}
