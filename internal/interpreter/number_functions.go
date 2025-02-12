package interpreter

import (
	"errors"
	"slices"
)

func makeOperatorFn(
	op string,
	leftAssoc bool) func([]ValueObject) (ValueObject, error) {

	binFn := makeBinOp(op)

	return func(numbers []ValueObject) (ValueObject, error) {
		var ret ValueObject
		var err error

		n := len(numbers)

		if n == 0 {
			return nil, errors.New("no arguments given for operator '" + op + "'")
		}
		if leftAssoc {
			ret = numbers[0]
			for _, number := range numbers[1:] {
				ret, err = binFn(ret, number)
				if err != nil {
					return nil, err
				}
			}
		} else {
			ret = numbers[n-1]
			for i := n - 2; i >= 0; i-- {
				ret, err = binFn(numbers[i], ret)
				if err != nil {
					return nil, err
				}
			}
		}

		return ret, nil
	}
}

func makeBinOp(op string) func(a, b ValueObject) (ValueObject, error) {
	return func(a, b ValueObject) (ValueObject, error) {
		return binOp(op, a, b)
	}
}

func binOp(op string, a, b ValueObject) (ValueObject, error) {
	numberType, err := getCommonNumberType(a.GetValueType(), b.GetValueType())
	if err != nil {
		return nil, err
	}
	numberA, err := convertNumber(a, numberType)
	if err != nil {
		return nil, err
	}
	numberB, err := convertNumber(b, numberType)
	if err != nil {
		return nil, err
	}
	switch numberType {
	case ValueInteger:
		intA := numberA.(*Integer)
		intB := numberB.(*Integer)
		switch op {
		case "+":
			return intA.Add(intB), nil
		case "-":
			return intA.Sub(intB), nil
		case "*":
			return intA.Mul(intB), nil
		case "/":
			return intA.Div(intB)
		case "%":
			return intA.Mod(intB), nil
		case "^":
			return intA.Pow(intB), nil
		default:
			return nil, errors.New("unknown operator: " + op)
		}
	case ValueRational:
		ratA := numberA.(*Rational)
		ratB := numberB.(*Rational)
		switch op {
		case "+":
			return ratA.Add(ratB), nil
		case "-":
			return ratA.Sub(ratB), nil
		case "*":
			return ratA.Mul(ratB), nil
		case "/":
			return ratA.Div(ratB)
		default:
			return nil, errors.New("unsupported operator for rational numbers: " + op)
		}
	case ValueReal:
		realA := numberA.(*Real)
		realB := numberB.(*Real)
		switch op {
		case "+":
			return realA.Add(realB), nil
		case "-":
			return realA.Sub(realB), nil
		case "*":
			return realA.Mul(realB), nil
		case "/":
			return realA.Div(realB)
		case "^":
			return realA.Pow(realB), nil
		default:
			return nil, errors.New("unsupported operator for real numbers: " + op)
		}
	default:
		return nil, errors.New("operator '" + op + "' is not supported")
	}
}

func convertNumber(value ValueObject, targetType ValueType) (ValueObject, error) {
	srcType := value.GetValueType()
	switch srcType {
	case ValueInteger:
		switch targetType {
		case ValueInteger:
			return value, nil
		case ValueRational:
			intVal := value.(*Integer)
			return intVal.ToRational(), nil
		case ValueReal:
			intVal := value.(*Integer)
			return intVal.ToReal(), nil
		default:
			return nil, errors.New("cannot convert integer")
		}
	case ValueRational:
		switch targetType {
		case ValueRational:
			return value, nil
		case ValueReal:
			rationalVal := value.(*Rational)
			return rationalVal.ToReal(), nil
		default:
			return nil, errors.New("cannot convert rational number")
		}
	case ValueReal:
		switch targetType {
		case ValueReal:
			return value, nil
		default:
			return nil, errors.New("cannot convert real number")
		}
	default:
		return nil, errors.New("cannot convert non-number")
	}
}

func compareNumbers(op string, numbers []ValueObject) (ValueObject, error) {
	commonType, err := getCommonTypeOfNums(numbers)
	if err != nil {
		return nil, err
	}

	var ints []*Integer
	var rationals []*Rational
	var reals []*Real

	for _, number := range numbers {
		convertedNum, errConv := convertNumber(number, commonType)
		if errConv != nil {
			return nil, errConv
		}
		switch commonType {
		case ValueInteger:
			ints = append(ints, convertedNum.(*Integer))
		case ValueRational:
			rationals = append(rationals, convertedNum.(*Rational))
		case ValueReal:
			reals = append(reals, convertedNum.(*Real))
		default:
			return nil, errors.New("invalid number type")
		}
	}

	switch commonType {
	case ValueInteger:
		return compare(op, ints)
	case ValueRational:
		return compare(op, rationals)
	case ValueReal:
		return compare(op, reals)
	default:
		return nil, errors.New("invalid number type")
	}
}

func getCommonTypeOfNums(numbers []ValueObject) (ValueType, error) {
	if len(numbers) == 0 {
		return ValueInvalid, errors.New("no numbers given")
	}

	var err error

	numberTypes := []ValueType{ValueInteger, ValueRational, ValueReal}
	ret := numbers[0].GetValueType()
	if !slices.Contains(numberTypes, ret) {
		return ValueInvalid, errors.New("value is not a number")
	}

	for _, number := range numbers[1:] {
		nextType := number.GetValueType()
		ret, err = getCommonNumberType(ret, nextType)
		if err != nil {
			return ValueInvalid, err
		}
	}

	return ret, nil
}

func getCommonNumberType(typeA, typeB ValueType) (ValueType, error) {
	numberTypes := []ValueType{ValueInteger, ValueRational, ValueReal}

	idxA := slices.Index(numberTypes, typeA)
	if idxA == -1 {
		return ValueInvalid, errors.New("not a number")
	}
	idxB := slices.Index(numberTypes, typeB)
	if idxB == -1 {
		return ValueInvalid, errors.New("not a number")
	}

	commonIdx := max(idxA, idxB)

	return numberTypes[commonIdx], nil
}

func integerDiv(objects []ValueObject) (ValueObject, error) {

	if len(objects) == 0 {
		return nil, errors.New("no integers given to // operator")
	}

	integerVals := make([]*Integer, len(objects))
	for i, object := range objects {
		if object.GetValueType() == ValueInteger {
			integerVals[i] = object.(*Integer)
		} else {
			return nil, errors.New("integer division operator // requires integers")
		}
	}

	result := integerVals[0].Value
	for _, intVal := range integerVals[1:] {
		result /= intVal.Value
	}

	return NewInteger(result), nil
}
