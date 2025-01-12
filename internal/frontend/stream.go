package frontend

import "errors"

type Stream[T interface{}] interface {
	Advance() (T, error)
	HasNext() bool
}

type BufferedStream[T interface{}] struct {
	stream Stream[T]
	buf    []T
}

func NewBufferedStream[T interface{}](stream Stream[T]) *BufferedStream[T] {
	return &BufferedStream[T]{
		stream: stream,
		buf:    make([]T, 0, 2),
	}
}

func (bufStream *BufferedStream[T]) Peek() (T, error) {
	var empty T
	nextItems := bufStream.PeekMany(1)
	if len(nextItems) == 0 {
		return empty, errors.New("no more items")
	}
	return nextItems[0], nil
}

func (bufStream *BufferedStream[T]) Advance() (T, error) {
	if len(bufStream.buf) > 0 {
		ret := bufStream.buf[0]
		bufStream.buf = bufStream.buf[1:]
		return ret, nil
	}

	return bufStream.stream.Advance()
}

func (bufStream *BufferedStream[T]) HasNext() bool {
	if len(bufStream.buf) > 0 {
		return true
	}

	return bufStream.stream.HasNext()
}

func (bufStream *BufferedStream[T]) PeekMany(n int) []T {
	for len(bufStream.buf) < n {
		item, err := bufStream.stream.Advance()
		if err != nil {
			break
		}
		bufStream.buf = append(bufStream.buf, item)
	}

	var size int
	if len(bufStream.buf) < n {
		size = len(bufStream.buf)
	} else {
		size = n
	}

	return bufStream.buf[:size]
}
