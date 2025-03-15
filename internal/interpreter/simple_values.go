package interpreter

import (
	"errors"
	"fmt"
	"go-boli/internal/frontend"
	"math"
	"strings"
)

type NilObject struct{}

var nilObject = NilObject{}

func GetNilObject() *NilObject {
	return &nilObject
}

func (n *NilObject) GetValueType() ValueType {
	return ValueNil
}

func (n *NilObject) String() string {
	return "nil"
}

func (n *NilObject) Car() (ValueObject, error) {
	return n, nil
}

func (n *NilObject) Cdr() (ValueObject, error) {
	return n, nil
}

func (n *NilObject) Take(int) (Sequence, error) {
	return n, nil
}

func (n *NilObject) TakeWhile(Callable) (Sequence, error) {
	return n, nil
}

func (n *NilObject) Filter(Callable) (Sequence, error) {
	return n, nil
}

func (n *NilObject) Map(Callable) (Sequence, error) {
	return n, nil
}

func (n *NilObject) Drop(int) (Sequence, error) {
	return n, nil
}

func (n *NilObject) DropWhile(Callable) (Sequence, error) {
	return n, nil
}

func (n *NilObject) Count() int {
	return 0
}

type Boolean struct {
	Value bool
}

func NewBoolean(v bool) *Boolean {
	return &Boolean{Value: v}
}

func (b *Boolean) GetValueType() ValueType {
	return ValueBoolean
}

func (b *Boolean) String() string {
	if b.Value {
		return "#true"
	} else {
		return "#false"
	}
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
	return fmt.Sprintf("%d", i.Value)
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

func (i *Integer) Equal(other *Integer) bool {
	return i.Value == other.Value
}

func (i *Integer) GreaterThan(other *Integer) bool {
	return i.Value > other.Value
}

func (i *Integer) GreaterThanOrEqual(other *Integer) bool {
	return i.Value >= other.Value
}

func (i *Integer) LessThan(other *Integer) bool {
	return i.Value < other.Value
}

func (i *Integer) LessThanOrEqual(other *Integer) bool {
	return i.Value <= other.Value
}

func (i *Integer) HashStr() string {
	return fmt.Sprintf("%d-%d", i.GetValueType(), i.Value)
}

func (i *Integer) IsEqual(other ValueObject) bool {
	return i.Value == other.(*Integer).Value
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

func (r *Rational) String() string {
	return fmt.Sprintf("%d/%d", r.Numerator, r.Denominator)
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

func (r *Rational) Equal(other *Rational) bool {
	return r.Numerator*other.Denominator == other.Numerator*r.Denominator
}

func (r *Rational) GreaterThan(other *Rational) bool {
	return r.Numerator*other.Denominator > other.Numerator*r.Denominator
}

func (r *Rational) GreaterThanOrEqual(other *Rational) bool {
	return r.Numerator*other.Denominator >= other.Numerator*r.Denominator
}

func (r *Rational) LessThan(other *Rational) bool {
	return r.Numerator*other.Denominator < other.Numerator*r.Denominator
}

func (r *Rational) LessThanOrEqual(other *Rational) bool {
	return r.Numerator*other.Denominator <= r.Numerator*r.Denominator
}

func (r *Rational) IsEqual(other ValueObject) bool {
	s := other.(*Rational)
	return r.Numerator == s.Numerator && r.Denominator == s.Denominator
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

func (r *Real) String() string {
	return strings.Replace(fmt.Sprintf("%f", r.Value), ".", ",", -1)
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

func (r *Real) Equal(other *Real) bool {
	const epsilon float64 = 1.0e-30
	return math.Abs(r.Value-other.Value) < epsilon
}

func (r *Real) GreaterThan(other *Real) bool {
	return r.Value > other.Value
}

func (r *Real) GreaterThanOrEqual(other *Real) bool {
	return r.Value >= other.Value
}

func (r *Real) LessThan(other *Real) bool {
	return r.Value < other.Value
}

func (r *Real) LessThanOrEqual(other *Real) bool {
	return r.Value <= other.Value
}

type Symbol struct {
	Value string
}

func NewSymbol(value string) *Symbol {
	return &Symbol{Value: value}
}

func (s *Symbol) GetValueType() ValueType {
	return ValueSymbol
}

func (s *Symbol) String() string {
	return s.Value
}

func (s *Symbol) HashStr() string {
	return fmt.Sprintf("%d-%s", s.GetValueType(), s.Value)
}

func (s *Symbol) IsEqual(other ValueObject) bool {
	return s.Value == other.(*Symbol).Value
}

type QuotedValue struct {
	token *frontend.Token
}

func NewQuotedValue(token *frontend.Token) *QuotedValue {
	return &QuotedValue{token: token}
}

func (q *QuotedValue) GetValueType() ValueType {
	return ValueQuoted
}

func (q *QuotedValue) String() string {
	return "'" + q.token.Lexeme
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
	if b < 0 {
		b = -b
		a = -a
	}
	if b == 1 {
		return NewInteger(a)
	}
	return NewRational(a, b)
}
