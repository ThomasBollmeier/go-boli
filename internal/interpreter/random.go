package interpreter

import (
	"fmt"
	"math/rand"
	"time"
)

var random *rand.Rand

func randomSeed(values []ValueObject) (ValueObject, error) {
	if len(values) > 1 {
		return nil, fmt.Errorf("random-seed expects 0 or 1 arguments, got %d", len(values))
	}

	if len(values) == 0 {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		intVal, ok := values[0].(*Integer)
		if !ok {
			return nil, fmt.Errorf("random-seed expects an integer as argument")
		}
		random = rand.New(rand.NewSource(int64(intVal.Value)))
	}

	return GetNilObject(), nil
}

func randomInteger(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("random-int expects 1 arguments, got %d", len(values))
	}
	intVal, ok := values[0].(*Integer)
	if !ok {
		return nil, fmt.Errorf("random-int expects an integer as argument")
	}

	if random == nil {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return NewInteger(random.Intn(intVal.Value)), nil
}

func randomReal(values []ValueObject) (ValueObject, error) {
	if len(values) > 0 {
		return nil, fmt.Errorf("random-real expects 0 arguments, got %d", len(values))
	}

	if random == nil {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return NewReal(random.Float64()), nil
}
