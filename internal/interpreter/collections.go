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
		ret := "(list "
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

	return fmt.Sprintf("(cons %s %s)", p.first, p.second)
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

func (p *Pair) Flatten() []ValueObject {
	ret := []ValueObject{p.first}
	curr := p.second
loop:
	for {
		switch curr.GetValueType() {
		case ValuePair:
			ret = append(ret, curr.(*Pair).first)
			curr = curr.(*Pair).second
		case ValueNil:
			break loop
		default:
			ret = append(ret, curr)
			break loop
		}
	}

	return ret
}

func (p *Pair) Take(n int) (Sequence, error) {
	var elements []ValueObject
	listElements := p.Flatten()
	for i, element := range listElements {
		if i >= n {
			break
		}
		elements = append(elements, element)
	}

	return elementsToSequence(elements), nil
}

func (p *Pair) TakeWhile(pred Callable) (Sequence, error) {
	var elements []ValueObject
	listElements := p.Flatten()
	for _, element := range listElements {
		predVal, err := Call(pred, []ValueObject{element})
		if err != nil {
			return nil, err
		}
		if !isTruthy(predVal) {
			break
		}
		elements = append(elements, element)
	}

	return elementsToSequence(elements), nil
}

func (p *Pair) Filter(pred Callable) (Sequence, error) {
	var elements []ValueObject
	listElements := p.Flatten()
	for _, element := range listElements {
		predVal, err := Call(pred, []ValueObject{element})
		if err != nil {
			return nil, err
		}
		if isTruthy(predVal) {
			elements = append(elements, element)
		}
	}
	return elementsToSequence(elements), nil
}

func (p *Pair) Map(fn Callable) (Sequence, error) {
	var elements []ValueObject
	listElements := p.Flatten()
	for _, element := range listElements {
		mappedVal, err := Call(fn, []ValueObject{element})
		if err != nil {
			return nil, err
		}
		elements = append(elements, mappedVal)
	}
	return elementsToSequence(elements), nil
}

func (p *Pair) Drop(n int) (Sequence, error) {
	var elements []ValueObject
	listElements := p.Flatten()
	for i, element := range listElements {
		if i < n {
			continue
		}
		elements = append(elements, element)
	}
	return elementsToSequence(elements), nil

}

func (p *Pair) DropWhile(pred Callable) (Sequence, error) {
	var elements []ValueObject
	listElements := p.Flatten()
	dropped := false
	for _, element := range listElements {
		if !dropped {
			predVal, err := Call(pred, []ValueObject{element})
			if err != nil {
				return nil, err
			}
			if isTruthy(predVal) {
				continue
			} else {
				dropped = true
				elements = append(elements, element)
			}
		} else {
			elements = append(elements, element)
		}
	}

	return elementsToSequence(elements), nil
}

func elementsToSequence(elements []ValueObject) Sequence {
	if len(elements) == 0 {
		return GetNilObject()
	}
	var ret *Pair
	numElements := len(elements)
	for i := numElements - 1; i >= 0; i-- {
		if ret == nil {
			ret = NewPair(elements[i], GetNilObject())
		} else {
			ret = NewPair(elements[i], ret)
		}
	}

	return ret
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

func (v *Vector) Take(n int) (Sequence, error) {
	var elements []ValueObject
	for i, element := range v.elements {
		if i >= n {
			break
		}
		elements = append(elements, element)
	}
	return NewVector(elements), nil
}

func (v *Vector) TakeWhile(pred Callable) (Sequence, error) {
	var elements []ValueObject
	for _, element := range v.elements {
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

func (v *Vector) Filter(pred Callable) (Sequence, error) {
	var elements []ValueObject
	for _, element := range v.elements {
		predVal, err := Call(pred, []ValueObject{element})
		if err != nil {
			return nil, err
		}
		if isTruthy(predVal) {
			elements = append(elements, element)
		}
	}
	return NewVector(elements), nil
}

func (v *Vector) Map(fn Callable) (Sequence, error) {
	var elements []ValueObject
	for _, element := range v.elements {
		mappedVal, err := Call(fn, []ValueObject{element})
		if err != nil {
			return nil, err
		}
		elements = append(elements, mappedVal)
	}
	return NewVector(elements), nil
}

func (v *Vector) Drop(n int) (Sequence, error) {
	var elements []ValueObject
	for i, element := range v.elements {
		if i < n {
			continue
		}
		elements = append(elements, element)
	}
	return NewVector(elements), nil

}

func (v *Vector) DropWhile(pred Callable) (Sequence, error) {
	var elements []ValueObject
	dropped := false
	for _, element := range v.elements {
		if !dropped {
			predVal, err := Call(pred, []ValueObject{element})
			if err != nil {
				return nil, err
			}
			if isTruthy(predVal) {
				continue
			} else {
				dropped = true
				elements = append(elements, element)
			}
		} else {
			elements = append(elements, element)
		}
	}

	return NewVector(elements), nil
}

func vector(values []ValueObject) (ValueObject, error) {
	return NewVector(values), nil
}

func isVec(values []ValueObject) (ValueObject, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("expected single arg, got %d", len(values))
	}

	return NewBoolean(values[0].GetValueType() == ValueVector), nil
}

func vectorGetRef(values []ValueObject) (ValueObject, error) {
	if len(values) != 2 {
		return nil, fmt.Errorf("expected two args, got %d", len(values))
	}
	var v *Vector
	first := values[0]
	switch first.GetValueType() {
	case ValueVector:
		v = first.(*Vector)
	default:
		return nil, fmt.Errorf("vector list as first arg")
	}
	second := values[1]
	if second.GetValueType() != ValueInteger {
		return nil, fmt.Errorf("expected integer as second arg")
	}
	idx := second.(*Integer).Value
	if idx < 0 {
		return nil, fmt.Errorf("expected positive integer as second arg")
	}

	if idx >= len(v.elements) {
		return nil, fmt.Errorf("invalid index")
	}

	return v.elements[idx], nil
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
	var ret ValueObject = GetNilObject()

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
