package interpreter

import "fmt"

type Stream interface {
	ValueObject
	Clonable[Stream]
	Next() (ValueObject, bool)
}

type ListStream struct {
	list ValueObject
}

func NewListStream(value ValueObject) (*ListStream, error) {
	if boolVal, _ := isList([]ValueObject{value}); !boolVal.(*Boolean).Value {
		return nil, fmt.Errorf("expected list, got %s", value.String())
	}
	return &ListStream{list: value}, nil
}

func (stream *ListStream) Clone() Stream {
	return &ListStream{
		list: stream.list,
	}
}

func (stream *ListStream) GetValueType() ValueType {
	return ValueStream
}

func (stream *ListStream) String() string {
	return "<list stream>"
}

func (stream *ListStream) Next() (ValueObject, bool) {
	pair, ok := stream.list.(*Pair)
	if !ok {
		return nil, false
	}
	stream.list = pair.second
	return pair.first, true
}

func isStream(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 value, got %d", len(objects))
	}

	return NewBoolean(objects[0].GetValueType() == ValueStream), nil
}

func listToStream(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("expected 1 value, got %d", len(objects))
	}

	return NewListStream(objects[0])
}

type FilteredStream struct {
	stream      Stream
	predicateFn Callable
}

func NewFilteredStream(stream Stream, predicateFn Callable) *FilteredStream {
	return &FilteredStream{
		stream,
		predicateFn,
	}
}

func (filtered *FilteredStream) GetValueType() ValueType {
	return ValueStream
}

func (filtered *FilteredStream) String() string {
	return "<filtered stream>"
}

func (filtered *FilteredStream) Clone() Stream {
	return NewFilteredStream(filtered.stream.Clone(), filtered.predicateFn)
}

func (filtered *FilteredStream) Next() (ValueObject, bool) {
	for {
		value, ok := filtered.stream.Next()
		if !ok {
			return nil, false
		}
		predVal, err := Call(filtered.predicateFn, []ValueObject{value})
		if err != nil {
			return nil, false
		}
		if isTruthy(predVal) {
			return value, true
		}
	}
}

func filter(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	predicateFn, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of filter must be a function")
	}
	stream, ok := objects[1].(Stream)
	if !ok {
		return nil, fmt.Errorf("second argument of filter must be a stream")
	}

	return NewFilteredStream(stream, predicateFn), nil
}

func take(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueInteger {
		return nil, fmt.Errorf("first argument of take must be an integer")
	}
	if objects[1].GetValueType() != ValueStream {
		return nil, fmt.Errorf("second argument of take must be a stream")
	}

	n := objects[0].(*Integer).Value
	stream := objects[1].(Stream).Clone()

	var elements []ValueObject
	for i := 0; i < n; i++ {
		element, ok := stream.Next()
		if !ok {
			break
		}
		elements = append(elements, element)
	}

	return NewVector(elements), nil
}
