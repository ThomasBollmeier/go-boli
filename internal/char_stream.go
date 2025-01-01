package internal

import (
	"errors"
)

type CharStream interface {
	Peek() (rune, error)
	Advance() (rune, error)
}

type BufferedCharStream struct {
	stream CharStream
	buf    []rune
}

func NewBufferedCharStream(stream CharStream) *BufferedCharStream {
	return &BufferedCharStream{
		stream: stream,
		buf:    make([]rune, 0, 2),
	}
}

func (bufCharStream *BufferedCharStream) Peek() (rune, error) {
	nextChars := bufCharStream.PeekMany(1)
	if len(nextChars) == 0 {
		return ' ', errors.New("no more characters")
	}
	return nextChars[0], nil
}

func (bufCharStream *BufferedCharStream) Advance() (rune, error) {
	if len(bufCharStream.buf) > 0 {
		ret := bufCharStream.buf[0]
		bufCharStream.buf = bufCharStream.buf[1:]
		return ret, nil
	}

	return bufCharStream.stream.Advance()
}

func (bufCharStream *BufferedCharStream) PeekMany(n int) []rune {
	for len(bufCharStream.buf) < n {
		ch, err := bufCharStream.stream.Advance()
		if err != nil {
			break
		}
		bufCharStream.buf = append(bufCharStream.buf, ch)
	}

	var size int
	if len(bufCharStream.buf) < n {
		size = len(bufCharStream.buf)
	} else {
		size = n
	}

	return bufCharStream.buf[:size]
}

type CharStreamString struct {
	chars []rune
	idx   int
}

func NewCharStreamString(text string) *CharStreamString {
	var chars []rune
	for _, char := range text {
		chars = append(chars, char)
	}
	return &CharStreamString{chars, 0}
}

func (stream *CharStreamString) Peek() (rune, error) {
	if stream.idx >= len(stream.chars) {
		return 0, errors.New("end of stream")
	}
	return stream.chars[stream.idx], nil
}

func (stream *CharStreamString) Advance() (rune, error) {
	if stream.idx >= len(stream.chars) {
		return 0, errors.New("end of stream")
	}
	ret := stream.chars[stream.idx]
	stream.idx++
	return ret, nil
}
