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
	return "<stream>"
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
