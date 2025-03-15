package interpreter

import (
	"errors"
	"fmt"
	"strings"
)

type Str struct {
	Value string
}

func NewStr(value string) *Str {
	return &Str{Value: value}
}

func (s *Str) GetValueType() ValueType {
	return ValueString
}

func (s *Str) String() string {
	return "\"" + s.Value + "\""
}

func (s *Str) HashStr() string {
	return fmt.Sprintf("%d-%s", s.GetValueType(), s.Value)
}

func (s *Str) IsEqual(other ValueObject) bool {
	return s.Value == other.(*Str).Value
}

func (s *Str) Car() (ValueObject, error) {
	if s.Value == "" {
		return nil, errors.New("car not possible for empty string")
	}
	chars := toChars(s.Value)
	return NewStr(chars[0]), nil
}

func (s *Str) Cdr() (ValueObject, error) {
	if s.Value == "" {
		return nil, errors.New("cdr not possible for empty string")
	}
	chars := toChars(s.Value)
	return NewStr(strings.Join(chars[1:], "")), nil
}

func (s *Str) Take(n int) (Sequence, error) {
	chars := toChars(s.Value)
	m := min(len(chars), n)
	return NewStr(strings.Join(chars[:m], "")), nil
}

func (s *Str) TakeWhile(pred Callable) (Sequence, error) {
	var charsTaken []string
	chars := toChars(s.Value)
	for _, ch := range chars {
		value, err := Call(pred, []ValueObject{NewStr(ch)})
		if err != nil {
			return nil, err
		}
		if !isTruthy(value) {
			break
		}
		charsTaken = append(charsTaken, ch)
	}
	return NewStr(strings.Join(charsTaken, "")), nil
}

func (s *Str) Filter(pred Callable) (Sequence, error) {
	var charsFiltered []string
	chars := toChars(s.Value)
	for _, ch := range chars {
		value, err := Call(pred, []ValueObject{NewStr(ch)})
		if err != nil {
			return nil, err
		}
		if isTruthy(value) {
			charsFiltered = append(charsFiltered, ch)
		}
	}
	return NewStr(strings.Join(charsFiltered, "")), nil
}

func (s *Str) Map(fn Callable) (Sequence, error) {
	var elements []ValueObject
	chars := toChars(s.Value)
	for _, ch := range chars {
		value, err := Call(fn, []ValueObject{NewStr(ch)})
		if err != nil {
			return nil, err
		}
		elements = append(elements, value)
	}
	return NewVector(elements), nil
}

func (s *Str) Drop(n int) (Sequence, error) {
	chars := toChars(s.Value)
	m := min(len(chars), n)
	return NewStr(strings.Join(chars[m:], "")), nil
}

func (s *Str) DropWhile(pred Callable) (Sequence, error) {
	chars := toChars(s.Value)
	index := -1
	for i, ch := range chars {
		value, err := Call(pred, []ValueObject{NewStr(ch)})
		if err != nil {
			return nil, err
		}
		if !isTruthy(value) {
			index = i
			break
		}
	}
	if index != -1 {
		return NewStr(strings.Join(chars[index:], "")), nil
	} else {
		return NewStr(""), nil
	}
}

func (s *Str) Count() int {
	return len(toChars(s.Value))
}

func toChars(s string) []string {
	ret := make([]string, 0)
	for _, c := range s {
		ret = append(ret, string(c))
	}
	return ret
}
