package main

import (
	"fmt"
	"go-boli/internal"
	"go-boli/internal/interpreter"
	"os"
)

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		internal.Repl()
	default:
		filePath := args[0]
		source, err := interpreter.NewFileSourceFactory().GetSource(filePath)
		if err != nil {
			panic(err)
		}
		value, err := interpreter.RunSource(source, args[1:])
		if err != nil {
			panic(err)
		}
		if value.GetValueType() != interpreter.ValueNil {
			fmt.Println(value)
		}
	}
}
