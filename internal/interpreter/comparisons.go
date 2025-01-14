package interpreter

import (
	"errors"
)

type Comparable[T interface{}] interface {
	Equal(T) bool
	GreaterThan(T) bool
	GreaterThanOrEqual(T) bool
	LessThan(T) bool
	LessThanOrEqual(T) bool
}

func compare[T Comparable[T]](op string, values []T) (ValueObject, error) {
	if len(values) < 2 {
		return nil, errors.New("not enough values to compare")
	}

	a := values[0]

	for i := 1; i < len(values); i++ {
		b := values[i]
		switch op {
		case "=":
			if !a.Equal(b) {
				return NewBoolean(false), nil
			}
		case ">":
			if !a.GreaterThan(b) {
				return NewBoolean(false), nil
			}
		case ">=":
			if !a.GreaterThanOrEqual(b) {
				return NewBoolean(false), nil
			}
		case "<":
			if !a.LessThan(b) {
				return NewBoolean(false), nil
			}
		case "<=":
			if !a.LessThanOrEqual(b) {
				return NewBoolean(false), nil
			}
		}
		a = b
	}

	return NewBoolean(true), nil
}
