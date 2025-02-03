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

type MappedStream struct {
	stream Stream
	mapFn  Callable
}

func NewMappedStream(stream Stream, mapFn Callable) *MappedStream {
	return &MappedStream{
		stream,
		mapFn,
	}
}

func (mapped *MappedStream) GetValueType() ValueType {
	return ValueStream
}

func (mapped *MappedStream) String() string {
	return "<mapped stream>"
}

func (mapped *MappedStream) Clone() Stream {
	return NewMappedStream(mapped.stream.Clone(), mapped.mapFn)
}

func (mapped *MappedStream) Next() (ValueObject, bool) {
	value, ok := mapped.stream.Next()
	if !ok {
		return nil, false
	}
	mappedVal, err := Call(mapped.mapFn, []ValueObject{value})
	if err != nil {
		return nil, false
	}
	return mappedVal, true
}

type DroppedStream struct {
	stream   Stream
	n        int
	dropDone bool
}

func NewDroppedStream(stream Stream, n int, dropDone bool) *DroppedStream {
	return &DroppedStream{
		stream,
		n,
		dropDone,
	}
}

func (dropped *DroppedStream) GetValueType() ValueType {
	return ValueStream
}

func (dropped *DroppedStream) String() string {
	return "<dropped stream>"
}

func (dropped *DroppedStream) Clone() Stream {
	return NewDroppedStream(dropped.stream.Clone(), dropped.n, dropped.dropDone)
}

func (dropped *DroppedStream) Next() (ValueObject, bool) {
	if !dropped.dropDone {
		for i := 0; i < dropped.n; i++ {
			_, ok := dropped.stream.Next()
			if !ok {
				return nil, false
			}
		}
		dropped.dropDone = true
	}

	return dropped.stream.Next()
}

type DroppedWhileStream struct {
	stream    Stream
	predicate Callable
	dropDone  bool
}

func NewDroppedWhileStream(stream Stream, predicate Callable, dropDone bool) *DroppedWhileStream {
	return &DroppedWhileStream{
		stream,
		predicate,
		dropDone,
	}
}

func (dropped *DroppedWhileStream) GetValueType() ValueType {
	return ValueStream
}

func (dropped *DroppedWhileStream) String() string {
	return "<dropped-while stream>"
}

func (dropped *DroppedWhileStream) Clone() Stream {
	return NewDroppedWhileStream(dropped.stream.Clone(),
		dropped.predicate,
		dropped.dropDone)
}

func (dropped *DroppedWhileStream) Next() (ValueObject, bool) {
	if !dropped.dropDone {
		for {
			value, ok := dropped.stream.Next()
			if !ok {
				return nil, false
			}
			predVal, err := Call(dropped.predicate, []ValueObject{value})
			if err != nil {
				return nil, false
			}
			if !isTruthy(predVal) {
				dropped.dropDone = true
				return value, true
			}
		}
	}

	return dropped.stream.Next()
}

// functions

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

func dropWhile(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}

	predicate, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("fist argument of drop-while must be a function, got %s", objects[0])
	}

	stream, ok := objects[1].(Stream)
	if !ok {
		return nil, fmt.Errorf("second argument of drop must be a stream, got %s", objects[1])
	}

	return NewDroppedWhileStream(stream, predicate, false), nil
}

func drop(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}

	intVal, ok := objects[0].(*Integer)
	if !ok {
		return nil, fmt.Errorf("fist argument of drop must be an integer, got %s", objects[0])
	}

	stream, ok := objects[1].(Stream)
	if !ok {
		return nil, fmt.Errorf("second argument of drop must be a stream, got %s", objects[1])
	}

	return NewDroppedStream(stream, intVal.Value, false), nil
}

func mapFunc(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}
	mapFn, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of map must be a function")
	}
	stream, ok := objects[1].(Stream)
	if !ok {
		return nil, fmt.Errorf("second argument of map must be a stream")
	}

	return NewMappedStream(stream, mapFn), nil
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

func takeWhile(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(objects))
	}

	predicate, ok := objects[0].(Callable)
	if !ok {
		return nil, fmt.Errorf("first argument of take-while must be a function")
	}
	if objects[1].GetValueType() != ValueStream {
		return nil, fmt.Errorf("second argument of take-while must be a stream")
	}

	stream := objects[1].(Stream).Clone()

	var elements []ValueObject
	for {
		element, ok := stream.Next()
		if !ok {
			break
		}
		predVal, err := Call(predicate, []ValueObject{element})
		if err != nil {
			return nil, err
		}
		if !isTruthy(predVal) {
			break
		}
		elements = append(elements, element)
	}

	return NewVector(elements), nil
}
