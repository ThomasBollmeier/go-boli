package interpreter

import "fmt"

type Sequence interface {
	ValueObject
	Take(n int) (Sequence, error)
	TakeWhile(pred Callable) (Sequence, error)
	Filter(pred Callable) (Sequence, error)
	Map(fn Callable) (Sequence, error)
	Drop(n int) (Sequence, error)
	DropWhile(pred Callable) (Sequence, error)
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
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	fn, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of map must be a function")
	}

	sequence, ok := objects[1].(Sequence)
	if !ok {
		return nil, fmt.Errorf("second argument of map must be a sequence")
	}

	return sequence.Map(fn)
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
