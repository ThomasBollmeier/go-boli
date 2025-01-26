package main

import (
	"fmt"
	"go-boli/internal"
	"go-boli/internal/interpreter"
	"os"
)

func main() {
	switch len(os.Args) {
	case 1:
		internal.Repl()
	case 2:
		value, err := interpreter.RunFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		if value.GetValueType() != interpreter.ValueNil {
			fmt.Println(value)
		}
	default:
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
	}
}
