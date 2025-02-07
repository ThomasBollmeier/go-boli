package interpreter

import "fmt"

type HashSet struct {
	entries map[string]Hashable
}

func NewHashSet(entries map[string]Hashable) *HashSet {
	return &HashSet{entries}
}

func (s *HashSet) GetValueType() ValueType {
	return ValueHashSet
}

func (s *HashSet) String() string {
	return "<set>"
}

func createSet(objects []ValueObject) (ValueObject, error) {
	entries := make(map[string]Hashable)

	for _, obj := range objects {
		hashable, ok := obj.(Hashable)
		if !ok {
			return nil, fmt.Errorf("expected hashable object")
		}
		entries[hashable.HashStr()] = hashable
	}

	return NewHashSet(entries), nil
}

func setLength(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("set-length expects 1 argument, got %d", len(objects))
	}

	hashSet, ok := objects[0].(*HashSet)
	if !ok {
		return nil, fmt.Errorf("expected a set as first argument of set-length")
	}

	return NewInteger(len(hashSet.entries)), nil
}

func setElements(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("set-elements expects 1 argument, got %d", len(objects))
	}

	hashSet, ok := objects[0].(*HashSet)
	if !ok {
		return nil, fmt.Errorf("expected a set as first argument of set-length")
	}

	var elements []ValueObject
	for _, entry := range hashSet.entries {
		elements = append(elements, entry)
	}

	return NewVector(elements), nil
}

func setContains(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("set-contains? expects 2 arguments, got %d", len(objects))
	}

	hashSet, ok := objects[0].(*HashSet)
	if !ok {
		return nil, fmt.Errorf("expected a set as first argument of set-contains?")
	}

	item, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of set-contains?")
	}

	_, ok = hashSet.entries[item.HashStr()]

	return NewBoolean(ok), nil
}

func setAddBang(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("set-add! expects 2 arguments, got %d", len(objects))
	}

	hashSet, ok := objects[0].(*HashSet)
	if !ok {
		return nil, fmt.Errorf("expected a set as first argument of set-add!")
	}

	item, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of set-add!")
	}

	hashSet.entries[item.HashStr()] = item

	return GetNilObject(), nil
}

func setRemoveBang(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("set-remove! expects  arguments, got %d", len(objects))
	}

	hashSet, ok := objects[0].(*HashSet)
	if !ok {
		return nil, fmt.Errorf("expected a set as first argument of set-remove!")
	}

	item, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of set-remove!")
	}

	delete(hashSet.entries, item.HashStr())

	return GetNilObject(), nil
}
