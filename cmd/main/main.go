package main

import (
	"fmt"
	"go-boli/internal"
	"go-boli/internal/interpreter"
	"os"
)

func main() {
	switch len(os.Args) {
	case 0:
		fmt.Printf("Usage: %s <file> [args...]\n", os.Args[0])
	case 1:
		internal.Repl()
	default:
		source, err := interpreter.NewFileSourceFactory().GetSource(os.Args[1])
		if err != nil {
			panic(err)
		}
		value, err := interpreter.RunSource(source, os.Args[2:])
		if err != nil {
			panic(err)
		}
		if value.GetValueType() != interpreter.ValueNil {
			fmt.Println(value)
		}
	}
}
