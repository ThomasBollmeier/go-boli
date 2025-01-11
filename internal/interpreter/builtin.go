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
			result, divErr := intA.Div(intB)
			return result, divErr
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
			return ratA.Div(ratB), nil
		default:
			return nil, errors.New("unsupported operator for rational numbers: " + op)
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
