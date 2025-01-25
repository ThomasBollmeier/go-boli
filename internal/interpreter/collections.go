package interpreter

import (
	"fmt"
)

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

	if b := p.IsList(); b.Value {
		var curr ValueObject
		ret := "'("
		curr = p
		first := true
	listLoop:
		for {
			switch curr.GetValueType() {
			case ValuePair:
				pair := curr.(*Pair)
				if !first {
					ret += " "
				} else {
					first = false
				}
				ret += fmt.Sprintf("%s", pair.first)
				curr = pair.second
			case ValueNil:
				ret += ")"
				break listLoop
			default:
				panic("unexpected value type")
			}
		}

		return ret
	}

	return fmt.Sprintf("'(%s . %s)", p.first, p.second)
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

type Vector struct {
	elements []ValueObject
}

func NewVector(elements []ValueObject) *Vector {
	return &Vector{
		elements: elements,
	}
}

func (v *Vector) GetValueType() ValueType {
	return ValueVector
}

func (v *Vector) String() string {
	elementsStr := ""
	for i, el := range v.elements {
		if i > 0 {
			elementsStr += " "
		}
		elementsStr += el.String()
	}
	return fmt.Sprintf("(vector %s)", elementsStr)
}

func (v *Vector) Append(element ValueObject) {
	v.elements = append(v.elements, element)
}

func (v *Vector) GetElements() []ValueObject {
	return v.elements
}

func car(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single arg, got %d", len(values))
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
		return nil, fmt.Errorf("expected single arg, got %d", len(values))
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

func isPair(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single arg, got %d", len(values))
	}

	value := values[0]
	switch value.GetValueType() {
	case ValuePair:
		return NewBoolean(true), nil
	default:

		return NewBoolean(false), nil
	}
}

func isList(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single arg, got %d", len(values))
	}

	value := values[0]
	switch value.GetValueType() {
	case ValuePair:
		return value.(*Pair).IsList(), nil
	case ValueNil:
		return NewBoolean(true), nil
	default:
		return NewBoolean(false), nil
	}
}

func cons(values []ValueObject) (ValueObject, error) {
	if len(values) != 2 {
		return nil, fmt.Errorf("expected two args, got %d", len(values))
	}

	a := values[0]
	b := values[1]
	return NewPair(a, b), nil
}

func list(values []ValueObject) (ValueObject, error) {
	ret := GetNilObject()

	for i := len(values) - 1; i > -1; i-- {
		ret = NewPair(values[i], ret)
	}

	return ret, nil
}

func listGetRef(values []ValueObject) (ValueObject, error) {
	if len(values) != 2 {
		return nil, fmt.Errorf("expected two args, got %d", len(values))
	}
	var pair *Pair
	lst := values[0]
	switch lst.GetValueType() {
	case ValuePair:
		pair = lst.(*Pair)
		if isAList := pair.IsList(); !isAList.Value {
			return nil, fmt.Errorf("expected non empty list as first arg")
		}
	default:
		return nil, fmt.Errorf("expected list as first arg")
	}
	idxVal := values[1]
	if idxVal.GetValueType() != ValueInteger {
		return nil, fmt.Errorf("expected integer as second arg")
	}
	idx := idxVal.(*Integer).Value
	if idx < 0 {
		return nil, fmt.Errorf("expected positive integer as second arg")
	}

	var curr ValueObject = pair
	for idx > 0 {
		p, ok := curr.(*Pair)
		if !ok {
			return nil, fmt.Errorf("invalid index")
		}
		curr = p.second
		idx--
	}

	return curr.(*Pair).first, nil
}
