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
	ValueHashSet
	ValueInvalid
)

type ValueObject interface {
	fmt.Stringer
	GetValueType() ValueType
}

type Clonable[T any] interface {
	Clone() T
}

type Equatable interface {
	IsEqual(other ValueObject) bool
}

func makeTypeCheckFunc(name string, valueType ValueType) func(objects []ValueObject) (ValueObject, error) {
	return func(objects []ValueObject) (ValueObject, error) {
		if len(objects) != 1 {
			return nil, fmt.Errorf("%s expects single arg, got %d", name, len(objects))
		}
		return NewBoolean(objects[0].GetValueType() == valueType), nil
	}
}

func allValuesEqual(values []ValueObject) (ValueObject, error) {
	for i, value := range values[:len(values)-1] {
		eq, err := valuesEqual(value, values[i+1])
		if err != nil {
			return nil, err
		}
		if !eq {
			return NewBoolean(false), nil
		}
	}

	return NewBoolean(true), nil
}

func valuesEqual(a, b ValueObject) (bool, error) {
	if a.GetValueType() != b.GetValueType() {
		return false, nil
	}
	eqA, ok := a.(Equatable)
	if !ok {
		return false, fmt.Errorf("'%s' cannot be compared for equality", a)
	}

	return eqA.IsEqual(b), nil
}
