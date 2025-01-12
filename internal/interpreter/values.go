package interpreter

import (
	"errors"
	"fmt"
	"math"
)

type ValueType int

const (
	ValueNil ValueType = iota
	ValueInteger
	ValueRational
	ValueReal
	ValueString
	ValueBuiltinFunc
	ValueInvalid
)

type ValueObject interface {
	GetValueType() ValueType
}

type Callable interface {
	Call(args []ValueObject) (ValueObject, error)
}

type NilObject struct{}

var nilObject = NilObject{}

func GetNilObject() ValueObject {
	return &nilObject
}

func (n *NilObject) GetValueType() ValueType {
	return ValueNil
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

func (i *Integer) ToRational() ValueObject {
	return NewRational(i.Value, 1)
}

func (i *Integer) ToReal() *Real {
	return NewReal(float64(i.Value))
}

func (i *Integer) Add(other *Integer) ValueObject {
	return NewInteger(i.Value + other.Value)
}

func (i *Integer) Sub(other *Integer) ValueObject {
	return NewInteger(i.Value - other.Value)
}

func (i *Integer) Mul(other *Integer) ValueObject {
	return NewInteger(i.Value * other.Value)
}

func (i *Integer) Div(other *Integer) (ValueObject, error) {
	if other.Value == 0 {
		return nil, errors.New("division by zero")
	}
	return newQuotient(i.Value, other.Value), nil
}

func (i *Integer) Mod(other *Integer) ValueObject {
	return NewInteger(i.Value % other.Value)
}

func (i *Integer) Pow(other *Integer) ValueObject {
	result := 1
	exp := other.Value
	if exp >= 0 {
		for exp > 0 {
			result *= i.Value
			exp -= 1
		}
		return NewInteger(result)
	} else {
		exp = -exp
		for exp > 0 {
			result *= i.Value
		}
		return NewRational(1, result)
	}
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

func (r *Rational) ToReal() *Real {
	return NewReal(float64(r.Numerator) / float64(r.Denominator))
}

func (r *Rational) Add(other *Rational) ValueObject {
	numerator := r.Numerator*other.Denominator + other.Numerator*r.Denominator
	denominator := r.Denominator * other.Denominator
	return newQuotient(numerator, denominator)
}

func (r *Rational) Sub(other *Rational) ValueObject {
	numerator := r.Numerator*other.Denominator - other.Numerator*r.Denominator
	denominator := r.Denominator * other.Denominator
	return newQuotient(numerator, denominator)
}

func (r *Rational) Mul(other *Rational) ValueObject {
	numerator := r.Numerator * other.Numerator
	denominator := r.Denominator * other.Denominator
	return newQuotient(numerator, denominator)
}

func (r *Rational) Div(other *Rational) (ValueObject, error) {
	numerator := r.Numerator * other.Denominator
	denominator := r.Denominator * other.Numerator
	if denominator == 0 {
		return nil, errors.New("division by zero")
	}
	return newQuotient(numerator, denominator), nil
}

func (r *Rational) String() string {
	return fmt.Sprintf("Rational(%d/%d)", r.Numerator, r.Denominator)
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

func (r *Real) Add(other *Real) ValueObject {
	return NewReal(r.Value + other.Value)
}

func (r *Real) Sub(other *Real) ValueObject {
	return NewReal(r.Value - other.Value)
}

func (r *Real) Mul(other *Real) ValueObject {
	return NewReal(r.Value * other.Value)
}

func (r *Real) Div(other *Real) (ValueObject, error) {
	if math.Abs(other.Value) < 1.0e-30 {
		return nil, errors.New("division by zero")
	}
	return NewReal(r.Value / other.Value), nil
}

func (r *Real) Pow(other *Real) ValueObject {
	return NewReal(math.Pow(r.Value, other.Value))
}

func (r *Real) String() string {
	return fmt.Sprintf("Real(%f)", r.Value)
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

func (b *BuiltinFunc) Call(args []ValueObject) (ValueObject, error) {
	return b.fn(args)
}

// Helper functions

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func shorten(a, b int) (int, int) {
	q := gcd(a, b)
	return a / q, b / q
}

func newQuotient(numerator int, denominator int) ValueObject {
	a, b := shorten(numerator, denominator)
	if b == 1 {
		return NewInteger(a)
	}
	return NewRational(a, b)
}
