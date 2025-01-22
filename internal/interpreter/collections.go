package interpreter

import "fmt"

type Pair struct {
	first  ValueObject
	second ValueObject
}

func NewPair(first ValueObject, second ValueObject) *Pair {
	return &Pair{
		first:  first,
		second: second,
	}
}

func (p *Pair) GetValueType() ValueType {
	return ValuePair
}

func (p *Pair) String() string {
	return fmt.Sprintf("( %s . %s )", p.first, p.second)
}

func (p *Pair) Car() ValueObject {
	return p.first
}

func (p *Pair) Cdr() ValueObject {
	return p.second
}

func (p *Pair) IsList() *Boolean {
	currValue := p.second
	for {
		switch currValue.GetValueType() {
		case ValuePair:
			currValue = currValue.(*Pair).second
		case ValueNil:
			return NewBoolean(true)
		default:
			return NewBoolean(false)
		}
	}
}

func car(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single value, got %d", len(values))
	}

	value := values[0]
	switch value.GetValueType() {
	case ValuePair:
		pair := value.(*Pair)
		return pair.Car(), nil
	default:
		return nil, fmt.Errorf("expected pair")
	}
}

func cdr(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single value, got %d", len(values))
	}

	value := values[0]
	switch value.GetValueType() {
	case ValuePair:
		pair := value.(*Pair)
		return pair.Cdr(), nil
	default:
		return nil, fmt.Errorf("expected pair")
	}
}

func isList(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single value, got %d", len(values))
	}

	value := values[0]
	switch value.GetValueType() {
	case ValuePair:
		return value.(*Pair).IsList(), nil
	default:
		return NewBoolean(false), nil
	}
}
