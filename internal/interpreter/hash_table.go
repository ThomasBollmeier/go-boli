package interpreter

import "fmt"

type Hashable interface {
	ValueObject
	HashStr() string
}

type HashEntry struct {
	key   Hashable
	value ValueObject
}

type HashTable struct {
	entries map[string]HashEntry
}

func NewHashTable() *HashTable {
	return &HashTable{
		entries: make(map[string]HashEntry),
	}
}

func (ht *HashTable) GetValueType() ValueType {
	return ValueHashTable
}

func (ht *HashTable) String() string {
	return "<hash table>"
}

func createHashTable(objects []ValueObject) (ValueObject, error) {
	if len(objects)%2 != 0 {
		return nil, fmt.Errorf("create-hash-table expects an even number of arguments, got %d", len(objects))
	}

	ret := NewHashTable()
	for i := 0; i < len(objects); i += 2 {
		key, ok := objects[i].(Hashable)
		if !ok {
			return nil, fmt.Errorf("expected a hashable value as argument %d of create-hash-table", i+1)
		}

		ret.entries[key.HashStr()] = HashEntry{
			key:   key,
			value: objects[i+1],
		}
	}

	return ret, nil
}

func hashLength(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("hash-length expects 1 argument, got %d", len(objects))
	}

	hashTable, ok := objects[0].(*HashTable)
	if !ok {
		return nil, fmt.Errorf("expected a hash table as first argument of hash-length")
	}

	return NewInteger(len(hashTable.entries)), nil
}

func hashKeys(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("hash-keys expects 1 argument, got %d", len(objects))
	}

	hashTable, ok := objects[0].(*HashTable)
	if !ok {
		return nil, fmt.Errorf("expected a hash table as first argument of hash-keys")
	}

	var keys []ValueObject
	for _, entry := range hashTable.entries {
		keys = append(keys, entry.key)
	}

	return NewVector(keys), nil
}

func hashContains(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("hash-contains expects 2 arguments, got %d", len(objects))
	}

	hashTable, ok := objects[0].(*HashTable)
	if !ok {
		return nil, fmt.Errorf("expected a hash table as first argument of hash-contains?")
	}

	key, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of hash-contains?")
	}

	_, ok = hashTable.entries[key.HashStr()]

	return NewBoolean(ok), nil
}

func hashGet(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("hash-get expects 2 arguments, got %d", len(objects))
	}

	hashTable, ok := objects[0].(*HashTable)
	if !ok {
		return nil, fmt.Errorf("expected a hash table as first argument of hash-get")
	}

	key, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of hash-get")
	}

	entry, ok := hashTable.entries[key.HashStr()]

	if ok {
		return entry.value, nil
	} else {
		return NewBoolean(false), nil
	}
}

func hashSetBang(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 3 {
		return nil, fmt.Errorf("hash-set! expects 3 arguments, got %d", len(objects))
	}

	hashTable, ok := objects[0].(*HashTable)
	if !ok {
		return nil, fmt.Errorf("expected a hash table as first argument of hash-set!")
	}

	key, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of hash-set!")
	}

	hashTable.entries[key.HashStr()] = HashEntry{
		key:   key,
		value: objects[2],
	}

	return GetNilObject(), nil
}

func hashRemoveBang(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("hash-remove! expects  arguments, got %d", len(objects))
	}

	hashTable, ok := objects[0].(*HashTable)
	if !ok {
		return nil, fmt.Errorf("expected a hash table as first argument of hash-remove!")
	}

	key, ok := objects[1].(Hashable)
	if !ok {
		return nil, fmt.Errorf("expected a hashable value as second argument of hash-remove!")
	}

	delete(hashTable.entries, key.HashStr())

	return GetNilObject(), nil
}
