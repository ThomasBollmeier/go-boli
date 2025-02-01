package interpreter

import (
	"errors"
	"fmt"
	"go-boli/internal/frontend"
)

type Callable interface {
	Call(args []ValueObject) (ValueObject, error)
}

func Call(callable Callable, arguments []ValueObject) (ValueObject, error) {
	var ret ValueObject
	var err error
	for {
		ret, err = callable.Call(arguments)
		if err != nil {
			return nil, err
		}
		switch ret.GetValueType() {
		case ValueTailCall:
			tailCall := ret.(*TailCall)
			callable = tailCall.callable
			arguments = tailCall.args
		default:
			return ret, nil
		}
	}
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

type LambdaVariant struct {
	params   []string
	varParam string
	body     *frontend.AST
	env      *Environment
}

func NewLambdaVariant(
	params []string,
	varParam string,
	body *frontend.AST,
	env *Environment) *LambdaVariant {
	return &LambdaVariant{
		params,
		varParam,
		body,
		env,
	}
}

func (l *LambdaVariant) GetMinMaxArgNum() (int, int) {
	minNumArgs := len(l.params)
	var maxNumArgs int
	if l.varParam == "" {
		maxNumArgs = minNumArgs
	} else {
		maxNumArgs = -1
	}

	return minNumArgs, maxNumArgs
}

func (l *LambdaVariant) OverlapsWith(other *LambdaVariant) bool {
	min1, max1 := l.GetMinMaxArgNum()
	min2, max2 := other.GetMinMaxArgNum()

	if max1 == -1 {
		if max2 == -1 {
			return true
		} else {
			return min1 <= max2
		}
	}

	if min1 >= min2 && (min1 <= max2 || max2 == -1) {
		return true
	}

	if max1 >= min2 && (max1 <= max2 || max2 == -1) {
		return true
	}

	return false
}

type LambdaFunc struct {
	name     string
	variants []*LambdaVariant
}

func NewLambdaFunc(
	name string,
	params []string,
	varParam string,
	body *frontend.AST,
	env *Environment) *LambdaFunc {

	return &LambdaFunc{
		name: name,
		variants: []*LambdaVariant{
			NewLambdaVariant(params, varParam, body, env),
		},
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

func (l *LambdaFunc) Merge(other *LambdaFunc) error {
	for _, variant := range other.variants {
		err := l.addVariant(variant)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LambdaFunc) addVariant(variant *LambdaVariant) error {
	for _, v := range l.variants {
		if v.OverlapsWith(variant) {
			return errors.New("new variant overlaps with existing lambda variant")
		}
	}

	l.variants = append(l.variants, variant)
	return nil
}

func (l *LambdaFunc) Call(args []ValueObject) (ValueObject, error) {
	numArgs := len(args)

	// Find matching variant:
	var matchingVariant *LambdaVariant
	for _, variant := range l.variants {
		minNumArgs, maxNumArgs := variant.GetMinMaxArgNum()
		if numArgs >= minNumArgs && (numArgs <= maxNumArgs || maxNumArgs == -1) {
			matchingVariant = variant
			break
		}
	}
	if matchingVariant == nil {
		return nil, errors.New("no variant with matching arity found")
	}

	interpreter := NewInterpreter(matchingVariant.env)
	interpreter.beginBlockScope()

	env := interpreter.env

	for i, param := range matchingVariant.params {
		env.Set(param, args[i], true)
	}

	numParams := len(matchingVariant.params)

	if matchingVariant.varParam != "" {
		varArgs := NewVector(nil)
		for i := numParams; i < numArgs; i++ {
			varArgs.Append(args[i])
		}
		env.Set(matchingVariant.varParam, varArgs, true)
	}

	ret, err := interpreter.evalBlock(matchingVariant.body)
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
