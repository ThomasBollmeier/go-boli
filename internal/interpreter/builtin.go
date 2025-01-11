package interpreter

import (
	"errors"
	"slices"
)

var fnAdd = NewBuiltinFunc("+", func(numbers []ValueObject) (ValueObject, error) {
	numberTypes := []ValueType{ValueInteger} //[]ValueType{ ValueInteger, ValueRational, ValueReal}
	var result int

	for _, number := range numbers {
		if !slices.Contains(numberTypes, number.GetValueType()) {
			return nil, errors.New("+ operator is only supported for numbers")
		}
		intVal := number.(*Integer)
		result += intVal.Value
	}

	return NewInteger(result), nil
})
