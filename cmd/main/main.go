package main

import (
	"fmt"
	"go-boli/internal/interpreter"
)

func main() {
	code := `(def answer 6)
	(* 7 answer)`
	value, err := interpreter.Run(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}
