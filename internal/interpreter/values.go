package interpreter

import "fmt"

type ValueType int

const (
	ValueNil ValueType = iota
	ValueBoolean
	ValueInteger
	ValueRational
	ValueReal
	ValueString
	ValueSymbol
	ValueQuoted
	ValueBuiltinFunc
	ValueLambda
	ValueTailCall
	ValuePair
	ValueVector
	ValueStream
	ValueStructType
	ValueStruct
	ValueHashTable
	ValueInvalid
)

type ValueObject interface {
	fmt.Stringer
	GetValueType() ValueType
}

type Clonable[T any] interface {
	Clone() T
}
