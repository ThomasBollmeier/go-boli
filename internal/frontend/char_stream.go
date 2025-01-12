package frontend

import (
	"errors"
)

type CharStream = Stream[rune]
type BufferedCharStream = BufferedStream[rune]

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

func (stream *CharStreamString) Advance() (rune, error) {
	if stream.idx >= len(stream.chars) {
		return 0, errors.New("end of stream")
	}
	ret := stream.chars[stream.idx]
	stream.idx++
	return ret, nil
}

func (stream *CharStreamString) HasNext() bool {
	return stream.idx < len(stream.chars)
}
