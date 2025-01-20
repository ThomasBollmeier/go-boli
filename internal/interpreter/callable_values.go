package interpreter

import (
	"errors"
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

func (b *BuiltinFunc) Call(args []ValueObject) (ValueObject, error) {
	return b.fn(args)
}

type LambdaFunc struct {
	name   string
	params []string
	body   *frontend.AST
	env    *Environment
}

func NewLambdaFunc(name string, params []string, body *frontend.AST, env *Environment) *LambdaFunc {
	return &LambdaFunc{
		name,
		params,
		body,
		env,
	}
}

func (l *LambdaFunc) GetValueType() ValueType {
	return ValueLambda
}

func (l *LambdaFunc) Call(args []ValueObject) (ValueObject, error) {
	if len(args) != len(l.params) {
		return nil, errors.New("invalid number of arguments")
	}
	interpreter := NewInterpreter(l.env)
	interpreter.beginBlockScope()

	env := interpreter.env
	for i, param := range l.params {
		env.Set(param, args[i])
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
