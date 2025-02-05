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

func makeTypeCheckFunc(name string, valueType ValueType) func(objects []ValueObject) (ValueObject, error) {
	return func(objects []ValueObject) (ValueObject, error) {
		if len(objects) != 1 {
			return nil, fmt.Errorf("%s expects single arg, got %d", name, len(objects))
		}
		return NewBoolean(objects[0].GetValueType() == valueType), nil
	}
}
