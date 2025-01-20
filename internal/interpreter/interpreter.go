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
	for _, op := range []string{"=", ">", ">=", "<", "<="} {
		ret.Set(op, NewBuiltinFunc(op, func(objects []ValueObject) (ValueObject, error) {
			return compareNumbers(op, objects)
		}))
	}

	return ret
}

func (interpreter *Interpreter) Eval(ast *frontend.AST) (ValueObject, error) {
	astType := ast.GetType()

	switch astType {
	case frontend.AstProgram:
		return interpreter.evalProgram(ast)
	case frontend.AstDefinition:
		return interpreter.evalDefinition(ast)
	case frontend.AstIfExpression:
		return interpreter.evalIfExpression(ast)
	case frontend.AstLambda:
		return interpreter.evalLambda(ast)
	case frontend.AstBlock:
		return interpreter.evalBlock(ast)
	case frontend.AstVariable:
		return interpreter.evalVariable(ast)
	case frontend.AstNil:
		return interpreter.evalNil()
	case frontend.AstBoolean:
		return interpreter.evalBoolean(ast)
	case frontend.AstInteger:
		return interpreter.evalInteger(ast)
	case frontend.AstRational:
		return interpreter.evalRational(ast)
	case frontend.AstReal:
		return interpreter.evalReal(ast)
	case frontend.AstString:
		return interpreter.evalString(ast)
	case frontend.AstOperator:
		return interpreter.evalOperator(ast)
	case frontend.AstComparisonOp:
		return interpreter.evalComparison(ast)
	case frontend.AstConjunction:
		return interpreter.evalConjunction(ast)
	case frontend.AstDisjunction:
		return interpreter.evalDisjunction(ast)
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

func (interpreter *Interpreter) evalComparison(op *frontend.AST) (ValueObject, error) {
	value, found := interpreter.env.Get(op.GetValue())
	if !found {
		return nil, errors.New(fmt.Sprintf("comparison operator '%s' is not defined", op.GetValue()))
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

func (interpreter *Interpreter) evalString(str *frontend.AST) (ValueObject, error) {
	text := str.GetValue()[1 : len(str.GetValue())-1]
	text = strings.Replace(text, "\\\"", "\"", -1)
	return NewStr(text), nil
}

func (interpreter *Interpreter) evalProgram(program *frontend.AST) (ValueObject, error) {
	var err error
	ret := GetNilObject()

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

func (interpreter *Interpreter) evalIfExpression(ifExpr *frontend.AST) (ValueObject, error) {
	children := ifExpr.GetChildren()
	condition, err := interpreter.Eval(children[0])
	if err != nil {
		return nil, err
	}

	if interpreter.isTruthy(condition) {
		return interpreter.Eval(children[1])
	} else {
		return interpreter.Eval(children[2])
	}
}

func (interpreter *Interpreter) evalLambda(lambda *frontend.AST) (ValueObject, error) {
	name := lambda.GetLexemes()
	children := lambda.GetChildren()
	var params []string
	for _, p := range children[0].GetChildren() {
		params = append(params, p.GetValue())
	}
	body := children[1]

	return NewLambdaFunc(name, params, body, interpreter.env), nil
}

func (interpreter *Interpreter) evalBlock(block *frontend.AST) (ValueObject, error) {
	var err error
	ret := GetNilObject()

	interpreter.beginBlockScope()

	for _, child := range block.GetChildren() {
		ret, err = interpreter.Eval(child)
		if err != nil {
			interpreter.endBlockScope()
			return nil, err
		}
	}

	interpreter.endBlockScope()

	return ret, nil
}

func (interpreter *Interpreter) evalVariable(variable *frontend.AST) (ValueObject, error) {
	varName := variable.GetValue()
	value, ok := interpreter.env.Get(varName)
	if !ok {
		return nil, errors.New(fmt.Sprintf("variable '%s' is not defined", varName))
	}
	return value, nil
}

func (interpreter *Interpreter) evalNil() (ValueObject, error) {
	return GetNilObject(), nil
}

func (interpreter *Interpreter) evalBoolean(ast *frontend.AST) (ValueObject, error) {
	switch ast.GetValue() {
	case "#true", "#t":
		return NewBoolean(true), nil
	case "#false", "#f":
		return NewBoolean(false), nil
	default:
		return nil, errors.New(fmt.Sprintf("boolean value '%s' is not defined", ast.GetValue()))
	}
}

func (interpreter *Interpreter) evalConjunction(conj *frontend.AST) (ValueObject, error) {
	var ret ValueObject = NewBoolean(true)
	var err error

	for _, child := range conj.GetChildren() {
		ret, err = interpreter.Eval(child)
		if err != nil {
			return nil, err
		}
		if !interpreter.isTruthy(ret) {
			return NewBoolean(false), nil
		}
	}

	return ret, nil
}

func (interpreter *Interpreter) evalDisjunction(disj *frontend.AST) (ValueObject, error) {
	var ret ValueObject = NewBoolean(true)
	var err error

	for _, child := range disj.GetChildren() {
		ret, err = interpreter.Eval(child)
		if err != nil {
			return nil, err
		}
		if interpreter.isTruthy(ret) {
			return ret, nil
		}
	}

	return NewBoolean(false), nil
}

func (interpreter *Interpreter) isTruthy(value ValueObject) bool {
	switch value.GetValueType() {
	case ValueNil:
		return false
	case ValueBoolean:
		boolValue, _ := value.(*Boolean)
		return boolValue.Value
	default:
		return true
	}
}

func (interpreter *Interpreter) beginBlockScope() {
	interpreter.env = NewEnvironment(interpreter.env)
}

func (interpreter *Interpreter) endBlockScope() {
	interpreter.env = interpreter.env.parent
}
