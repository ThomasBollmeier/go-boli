package main

import (
	"bitbucket.org/drbolle/go-boli/internal/interpreter"
	"fmt"
)

func main() {
	code := `(+ 83/2 2/4)`
	value, err := interpreter.Run(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}
