package interpreter

import (
	"fmt"
)

type Sequence interface {
	ValueObject
	Car() (ValueObject, error)
	Cdr() (ValueObject, error)
	Take(n int) (Sequence, error)
	TakeWhile(pred Callable) (Sequence, error)
	Filter(pred Callable) (Sequence, error)
	Map(fn Callable, sequences []Sequence) (Sequence, error)
	Drop(n int) (Sequence, error)
	DropWhile(pred Callable) (Sequence, error)
}

type LimitedSequence interface {
	Sequence
	Count() int
}

func car(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected single arg, got %d", len(objects))
	}

	sequence, ok := objects[0].(Sequence)
	if !ok {
		return nil, fmt.Errorf("car requires a sequence")
	}
	return sequence.Car()
}

func cdr(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected single arg, got %d", len(objects))
	}
	sequence, ok := objects[0].(Sequence)
	if !ok {
		return nil, fmt.Errorf("cdr requires a sequence")
	}
	return sequence.Cdr()
}

func take(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	intVal, ok := objects[0].(*Integer)
	if !ok {
		return nil, fmt.Errorf("first argument of take must be an integer")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of take must be a sequence")
	}

	return sequence.Take(intVal.Value)
}

func takeWhile(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	pred, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of take-while must be a function")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of take-while must be a sequence")
	}

	return sequence.TakeWhile(pred)
}

func filter(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	pred, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of filter must be a function")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of filter must be a sequence")
	}

	return sequence.Filter(pred)
}

func mapFunc(objects []ValueObject) (ValueObject, error) {
	if len(objects) < 2 {
		return nil, fmt.Errorf("expected at least two values, got %d", len(objects))
	}
	fn, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of map must be a function")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of map must be a sequence")
	}

	var otherSequences []Sequence

	for _, object := range objects[2:] {
		otherSequence, okOther := object.(Sequence)
		if !okOther {
			return nil, fmt.Errorf("arguments of map after first must be sequencea")
		}
		otherSequences = append(otherSequences, otherSequence)
	}

	return sequence.Map(fn, otherSequences)
}

func drop(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	intVal, ok := objects[0].(*Integer)
	if !ok {
		return nil, fmt.Errorf("first argument of drop must be an integer")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of drop must be a sequence")
	}

	return sequence.Drop(intVal.Value)
}

func dropWhile(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	pred, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of drop-while must be a function")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of drop-while must be a sequence")
	}

	return sequence.DropWhile(pred)
}

func count(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(objects))
	}
	lseq, ok := objects[0].(LimitedSequence)
	if !ok {
		return nil, fmt.Errorf("first argument of count must be a limited sequence")
	}

	return NewInteger(lseq.Count()), nil
}

func isEmpty(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(objects))
	}
	lseq, ok := objects[0].(LimitedSequence)
	if !ok {
		return nil, fmt.Errorf("first argument of empty? must be a limited sequence")
	}

	return NewBoolean(lseq.Count() == 0), nil
}

func splitFirstElements(sequences []Sequence) ([]ValueObject, []Sequence, error) {
	var elements []ValueObject
	var nextSequences []Sequence

	for _, sequence := range sequences {
		element, err := sequence.Car()
		if err != nil {
			return nil, nil, err
		}
		elements = append(elements, element)
		nextSequence, err := sequence.Cdr()
		if err != nil {
			return nil, nil, err
		}
		nextSequences = append(nextSequences, nextSequence.(Sequence))
	}

	return elements, nextSequences, nil
}
