package interpreter

import (
	"errors"
	"fmt"
	"go-boli/internal/frontend"
)

type Callable interface {
	Call(args []ValueObject) (ValueObject, error)
}

type BuiltinFunc struct {
	name string
	fn   func([]ValueObject) (ValueObject, error)
}

func NewBuiltinFunc(name string, fn func([]ValueObject) (ValueObject, error)) *BuiltinFunc {
	return &BuiltinFunc{name, fn}
}

func (b *BuiltinFunc) GetValueType() ValueType {
	return ValueBuiltinFunc
}

func (b *BuiltinFunc) String() string {
	return fmt.Sprintf("<builtin function %s>", b.name)
}

func (b *BuiltinFunc) Call(args []ValueObject) (ValueObject, error) {
	return b.fn(args)
}

type LambdaFunc struct {
	name     string
	params   []string
	varParam string
	body     *frontend.AST
	env      *Environment
}

func NewLambdaFunc(
	name string,
	params []string,
	varParam string,
	body *frontend.AST,
	env *Environment) *LambdaFunc {

	return &LambdaFunc{
		name,
		params,
		varParam,
		body,
		env,
	}
}

func (l *LambdaFunc) GetValueType() ValueType {
	return ValueLambda
}

func (l *LambdaFunc) String() string {
	if len(l.name) == 0 {
		return "<lambda function>"
	} else {
		return fmt.Sprintf("<lambda function %s>", l.name)
	}
}

func (l *LambdaFunc) Call(args []ValueObject) (ValueObject, error) {
	numArgs := len(args)
	numParams := len(l.params)

	if numArgs < numParams {
		return nil, errors.New("too few arguments given")
	} else if numArgs > numParams && l.varParam == "" {
		return nil, errors.New("too many arguments given")
	}

	interpreter := NewInterpreter(l.env)
	interpreter.beginBlockScope()

	env := interpreter.env

	for i, param := range l.params {
		env.Set(param, args[i])
	}

	if l.varParam != "" {
		varArgs := NewVector(nil)
		for i := numParams; i < numArgs; i++ {
			varArgs.Append(args[i])
		}
		env.Set(l.varParam, varArgs)
	}

	ret, err := interpreter.evalBlock(l.body)
	if err != nil {
		interpreter.endBlockScope()
		return nil, err
	}

	interpreter.endBlockScope()

	return ret, nil
}

type TailCall struct {
	callable Callable
	args     []ValueObject
}

func NewTailCall(callable Callable, args []ValueObject) *TailCall {
	return &TailCall{callable, args}
}

func (tc *TailCall) GetValueType() ValueType {
	return ValueTailCall
}

func (tc *TailCall) String() string {
	return "<tail call>"
}
