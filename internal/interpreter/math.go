package interpreter

import (
	"fmt"
	"math"
)

func NewMath() *BuiltInModule {
	result := NewBuiltInModule()
	result.AddValue("sqrt", makeSingleRealArgFn("sqrt", math.Sqrt))
	result.AddValue("exp", makeSingleRealArgFn("exp", math.Exp))
	return result
}

func makeSingleRealArgFn(name string, fn func(float64) float64) *BuiltinFunc {
	return NewBuiltinFunc(name, func(args []ValueObject) (ValueObject, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("%s requires exactly one argument", name)
		}

		arg := args[0]
		if arg.GetValueType() != ValueReal {
			return nil, fmt.Errorf("%s requires a real argument", name)
		}

		realVal, _ := arg.(*Real)

		return NewReal(fn(realVal.Value)), nil
	})
}
