package interpreter

import (
	"fmt"
	"strings"
)

type StructureType struct {
	name   string
	fields []string
}

func NewStructureType(name string, fields []string) *StructureType {
	return &StructureType{
		name:   name,
		fields: fields,
	}
}

func (structType *StructureType) GetValueType() ValueType {
	return ValueStructType
}

func (structType *StructureType) String() string {
	ret := "(struct " + structType.name + " ("
	ret += strings.Join(structType.fields, " ")
	ret += "))"
	return ret
}

func (structType *StructureType) createConstructor() *BuiltinFunc {
	name := "create-" + structType.name

	return NewBuiltinFunc(name, func(args []ValueObject) (ValueObject, error) {
		if len(args) != len(structType.fields) {
			return nil, fmt.Errorf("create method requires %v arguments, got %v",
				len(structType.fields), len(args))
		}

		values := make(map[string]ValueObject)
		for i, field := range structType.fields {
			values[field] = args[i]
		}

		return NewStructure(structType, values), nil
	})
}

type Structure struct {
	structType *StructureType
	values     map[string]ValueObject
}

func NewStructure(structType *StructureType, values map[string]ValueObject) *Structure {
	return &Structure{
		structType: structType,
		values:     values,
	}
}

func (structure *Structure) GetValueType() ValueType {
	return ValueStruct
}

func (structure *Structure) String() string {
	ret := fmt.Sprintf("(create-%s ", structure.structType.name)
	for i, field := range structure.structType.fields {
		if i > 0 {
			ret += " "
		}
		value := structure.values[field]
		ret += value.String()
	}
	ret += ")"
	return ret
}
