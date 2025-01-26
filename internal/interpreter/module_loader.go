package interpreter

import (
	"errors"
	"os"
)

const MODULE_SEPARATOR string = "::"

func LoadModule(moduleName, alias string) (*Environment, error) {
	return nil, errors.New("not implemented")
}

func LoadFile(filePath string) (*Environment, error) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter(nil)
	_, err = interpreter.Run(string(code))
	if err != nil {
		return nil, err
	}

	return interpreter.env, nil
}

func RunFile(filePath string) (ValueObject, error) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter(nil)
	value, err := interpreter.Run(string(code))
	if err != nil {
		return nil, err
	}

	return value, nil
}
