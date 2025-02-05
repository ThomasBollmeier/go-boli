package interpreter

import (
	"errors"
	"fmt"
)

type Stream interface {
	fmt.Stringer
	Clonable[Stream]
	Next() (ValueObject, bool)
}

type IterStream struct {
	current ValueObject
	done    bool
	nextFn  Callable
}

func NewIterStream(
	current ValueObject,
	done bool,
	nextFn Callable) *IterStream {
	return &IterStream{
		current,
		done,
		nextFn,
	}
}

func (it *IterStream) Clone() Stream {
	return NewIterStream(it.current, it.done, it.nextFn)
}

func (it *IterStream) String() string {
	return "<iterator stream>"
}

func (it *IterStream) Next() (ValueObject, bool) {
	if it.done {
		return nil, false
	}
	var err error
	ret := it.current
	it.current, err = Call(it.nextFn, []ValueObject{it.current})
	if err != nil || it.current.GetValueType() == ValueNil {
		it.done = true
	}
	return ret, true
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

type StreamSeq struct {
	stream Stream
}

func NewStreamSeq(stream Stream) *StreamSeq {
	return &StreamSeq{stream}
}

func (seq *StreamSeq) GetValueType() ValueType {
	return ValueStream
}

func (seq *StreamSeq) String() string {
	return seq.stream.String()
}

func (seq *StreamSeq) Car() (ValueObject, error) {
	first, err := seq.Take(1)
	if err != nil {
		return nil, err
	}
	firstVec, _ := first.(*Vector)
	if len(firstVec.elements) == 0 {
		return nil, errors.New("stream has no elements")
	}
	return firstVec.elements[0], nil
}

func (seq *StreamSeq) Cdr() (ValueObject, error) {
	return seq.Drop(1)
}

func (seq *StreamSeq) Take(n int) (Sequence, error) {
	var elements []ValueObject
	stream := seq.stream.Clone()
	for i := 0; i < n; i++ {
		element, ok := stream.Next()
		if !ok {
			break
		}
		elements = append(elements, element)
	}

	return NewVector(elements), nil
}

func (seq *StreamSeq) TakeWhile(pred Callable) (Sequence, error) {
	stream := seq.stream.Clone()
	var elements []ValueObject
	for {
		element, ok := stream.Next()
		if !ok {
			break
		}
		predVal, err := Call(pred, []ValueObject{element})
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

func (seq *StreamSeq) Filter(pred Callable) (Sequence, error) {
	return NewStreamSeq(NewFilteredStream(seq.stream, pred)), nil
}

func (seq *StreamSeq) Map(fn Callable) (Sequence, error) {
	return NewStreamSeq(NewMappedStream(seq.stream, fn)), nil
}

func (seq *StreamSeq) Drop(n int) (Sequence, error) {
	return NewStreamSeq(NewDroppedStream(seq.stream, n, false)), nil
}

func (seq *StreamSeq) DropWhile(pred Callable) (Sequence, error) {
	return NewStreamSeq(NewDroppedWhileStream(seq.stream, pred, false)), nil
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
	listStream, err := NewListStream(objects[0])
	if err != nil {
		return nil, err
	}

	return NewStreamSeq(listStream), nil
}

func iterator(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("iterator expects 2 arguments, got %d", len(objects))
	}
	nextFn, ok := objects[1].(Callable)
	if !ok {
		return nil, fmt.Errorf("expected a function as second argument of iterator")
	}

	return NewStreamSeq(NewIterStream(objects[0], false, nextFn)), nil
}
