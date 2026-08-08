package interpreter

import (
	"errors"
	"fmt"
	"strings"
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
		var ret strings.Builder
		ret.WriteString("(list ")
		curr = p
		first := true
	listLoop:
		for {
			switch curr.GetValueType() {
			case ValuePair:
				pair := curr.(*Pair)
				if !first {
					ret.WriteString(" ")
				} else {
					first = false
				}
				ret.WriteString(fmt.Sprintf("%s", pair.first))
				curr = pair.second
			case ValueNil:
				ret.WriteString(")")
				break listLoop
			default:
				panic("unexpected value type")
			}
		}

		return ret.String()
	}

	return fmt.Sprintf("(cons %s %s)", p.first, p.second)
}

func (p *Pair) Car() (ValueObject, error) {
	return p.first, nil
}

func (p *Pair) Cdr() (ValueObject, error) {
	return p.second, nil
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

func (p *Pair) Map(fn Callable, otherSequences []Sequence) (Sequence, error) {
	var elements []ValueObject
	var args []ValueObject
	var err error

	listElements := p.Flatten()
	sequences := otherSequences
	for _, element := range listElements {
		args, sequences, err = splitFirstElements(sequences)
		if err != nil {
			break
		}
		args = append([]ValueObject{element}, args...)
		mappedVal, errCall := Call(fn, args)
		if errCall != nil {
			return nil, errCall
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

func (p *Pair) Count() int {
	ret := 0
	var curr ValueObject = p
loop:
	for {
		switch curr.GetValueType() {
		case ValuePair:
			ret++
			pair := curr.(*Pair)
			curr = pair.second
		case ValueNil:
			break loop
		default:
			ret++
			break loop
		}
	}

	return ret
}

func (p *Pair) IsEqual(other ValueObject) bool {
	if other.GetValueType() != ValuePair {
		return false
	}
	op := other.(*Pair)
	ok, err := valuesEqual(p.first, op.first)
	if !ok || err != nil {
		return false
	}
	ok, err = valuesEqual(p.second, op.second)
	if !ok || err != nil {
		return false
	}
	return true
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

func (v *Vector) Car() (ValueObject, error) {
	if len(v.elements) == 0 {
		return nil, errors.New("cannot car from empty vector")
	}
	return v.elements[0], nil
}

func (v *Vector) Cdr() (ValueObject, error) {
	if len(v.elements) == 0 {
		return nil, errors.New("cannot cdr from empty vector")
	}
	return NewVector(v.elements[1:]), nil
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

func (v *Vector) Map(fn Callable, otherSequences []Sequence) (Sequence, error) {
	var elements []ValueObject
	var args []ValueObject
	var err error
	sequences := otherSequences

	for _, element := range v.elements {
		args, sequences, err = splitFirstElements(sequences)
		if err != nil {
			break
		}
		args = append([]ValueObject{element}, args...)
		mappedVal, errCall := Call(fn, args)
		if errCall != nil {
			return nil, errCall
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

func (v *Vector) Count() int {
	return len(v.elements)
}

func (v *Vector) IsEqual(other ValueObject) bool {
	if other.GetValueType() != ValueVector {
		return false
	}
	ov := other.(*Vector)
	if len(v.elements) != len(ov.elements) {
		return false
	}
	for i := 0; i < len(v.elements); i++ {
		elem := v.elements[i]
		otherElem := ov.elements[i]
		ok, err := valuesEqual(elem, otherElem)
		if !ok || err != nil {
			return false
		}
	}
	return true
}

func vector(values []ValueObject) (ValueObject, error) {
	return NewVector(values), nil
}

func vectorRef(values []ValueObject) (ValueObject, error) {
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

func vectorSetBang(values []ValueObject) (ValueObject, error) {
	if len(values) != 3 {
		return nil, fmt.Errorf("expected three args, got %d", len(values))
	}
	var v *Vector
	first := values[0]
	switch first.GetValueType() {
	case ValueVector:
		v = first.(*Vector)
	default:
		return nil, fmt.Errorf("expected vector as first arg")
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

	v.elements[idx] = values[2]

	return GetNilObject(), nil
}

func vectorAppend(values []ValueObject) (ValueObject, error) {
	if len(values) != 2 {
		return nil, fmt.Errorf("expected two args, got %d", len(values))
	}
	var v *Vector
	first := values[0]
	switch first.GetValueType() {
	case ValueVector:
		v = first.(*Vector)
	default:
		return nil, fmt.Errorf("expected vector as first arg")
	}

	v.elements = append(v.elements, values[1])

	return GetNilObject(), nil
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

func listRef(values []ValueObject) (ValueObject, error) {
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
	for {
		p, ok := curr.(*Pair)
		if !ok {
			return nil, fmt.Errorf("invalid index")
		}
		if idx == 0 {
			return p.first, nil
		}
		curr = p.second
		idx--
	}
}
