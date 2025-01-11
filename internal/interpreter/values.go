package interpreter

import "fmt"

type ValueType int

const (
	ValueInteger ValueType = iota
	ValueRational
	ValueReal
	ValueString
	ValueBuiltinFunc
)

type ValueObject interface {
	GetValueType() ValueType
}

type Callable interface {
	Call(params []ValueObject) (ValueObject, error)
}

type Integer struct {
	Value int
}

func NewInteger(value int) *Integer {
	return &Integer{Value: value}
}

func (i *Integer) GetValueType() ValueType {
	return ValueInteger
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%d)", i.Value)
}

type Rational struct {
	Numerator   int
	Denominator int
}

func NewRational(numerator, denominator int) *Rational {
	return &Rational{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

func (r *Rational) GetValueType() ValueType {
	return ValueRational
}

type Real struct {
	Value float64
}

func NewReal(value float64) *Real {
	return &Real{Value: value}
}

func (r *Real) GetValueType() ValueType {
	return ValueReal
}

type Str struct {
	Value string
}

func NewString(value string) *Str {
	return &Str{Value: value}
}

func (s *Str) GetValueType() ValueType {
	return ValueString
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

func (b *BuiltinFunc) Call(params []ValueObject) (ValueObject, error) {
	return b.fn(params)
}
