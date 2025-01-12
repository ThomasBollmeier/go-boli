package interpreter

import (
	"errors"
	"fmt"
	"go-boli/internal/frontend"
	"strconv"
	"strings"
)

func Run(code string) (ValueObject, error) {
	parser := frontend.NewParser()
	ast, err := parser.Parse(frontend.NewCharStreamString(code))
	if err != nil {
		return nil, err
	}
	interpreter := NewInterpreter(nil)
	return interpreter.Eval(ast)
}

type Interpreter struct {
	env *Environment
}

func NewInterpreter(env *Environment) *Interpreter {
	if env == nil {
		return &Interpreter{
			env: newGlobalEnv(),
		}
	}
	return &Interpreter{
		env: env,
	}
}

func newGlobalEnv() *Environment {
	ret := NewEnvironment(nil)
	for _, op := range []string{"+", "-", "*", "/", "%"} {
		ret.Set(op, NewBuiltinFunc(op, makeOperatorFn(op, true)))
	}
	ret.Set("^", NewBuiltinFunc("^", makeOperatorFn("^", false)))

	return ret
}

func (interpreter *Interpreter) Eval(ast *frontend.AST) (ValueObject, error) {
	astType := ast.GetType()

	switch astType {
	case frontend.AstProgram:
		return interpreter.evalProgram(ast)
	case frontend.AstDefinition:
		return interpreter.evalDefinition(ast)
	case frontend.AstVariable:
		return interpreter.evalVariable(ast)
	case frontend.AstInteger:
		return interpreter.evalInteger(ast)
	case frontend.AstRational:
		return interpreter.evalRational(ast)
	case frontend.AstReal:
		return interpreter.evalReal(ast)
	case frontend.AstOperator:
		return interpreter.evalOperator(ast)
	case frontend.AstCall:
		return interpreter.evalCall(ast)
	default:
		return nil, errors.New("not implemented")
	}
}

func (interpreter *Interpreter) evalInteger(ast *frontend.AST) (ValueObject, error) {
	value, _ := strconv.ParseInt(ast.GetValue(), 10, 0)
	return NewInteger(int(value)), nil
}

func (interpreter *Interpreter) evalCall(call *frontend.AST) (ValueObject, error) {
	children := call.GetChildren()
	calleeAst := children[0]
	callee, err := interpreter.Eval(calleeAst)
	if err != nil {
		return nil, err
	}
	callable, ok := callee.(Callable)
	if !ok {
		return nil, errors.New("callee is not a callable")
	}
	var arguments []ValueObject
	for _, child := range children[1:] {
		argument, childErr := interpreter.Eval(child)
		if childErr != nil {
			return nil, childErr
		}
		arguments = append(arguments, argument)
	}

	return callable.Call(arguments)
}

func (interpreter *Interpreter) evalOperator(op *frontend.AST) (ValueObject, error) {
	value, found := interpreter.env.Get(op.GetValue())
	if !found {
		return nil, errors.New(fmt.Sprintf("operator '%s' is not defined", op.GetValue()))
	}
	return value, nil
}

func (interpreter *Interpreter) evalRational(rational *frontend.AST) (ValueObject, error) {
	parts := strings.Split(rational.GetValue(), "/")
	if len(parts) != 2 {
		return nil, errors.New("rational must be of the form 'a/b'")
	}
	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	return newQuotient(numerator, denominator), nil
}

func (interpreter *Interpreter) evalReal(real *frontend.AST) (ValueObject, error) {
	realString := strings.Replace(real.GetValue(), ",", ".", -1)
	realValue, err := strconv.ParseFloat(realString, 0)
	if err != nil {
		return nil, err
	}
	return NewReal(realValue), nil
}

func (interpreter *Interpreter) evalProgram(program *frontend.AST) (ValueObject, error) {
	var ret ValueObject
	var err error

	for _, child := range program.GetChildren() {
		ret, err = interpreter.Eval(child)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (interpreter *Interpreter) evalDefinition(def *frontend.AST) (ValueObject, error) {
	name := def.GetValue()
	valueAst := def.GetChildren()[0]
	value, err := interpreter.Eval(valueAst)
	if err != nil {
		return nil, err
	}
	interpreter.env.Set(name, value)

	return GetNilObject(), nil
}

func (interpreter *Interpreter) evalVariable(variable *frontend.AST) (ValueObject, error) {
	varName := variable.GetValue()
	value, ok := interpreter.env.Get(varName)
	if !ok {
		return nil, errors.New(fmt.Sprintf("variable '%s' is not defined", varName))
	}
	return value, nil
}
