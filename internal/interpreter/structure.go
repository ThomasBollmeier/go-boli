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

func (structType *StructureType) createTypeChecker() *BuiltinFunc {
	name := structType.name + "?"
	return NewBuiltinFunc(name, func(args []ValueObject) (ValueObject, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("create method requires a single arguments, got %d", len(args))
		}
		arg0 := args[0]
		if arg0.GetValueType() == ValueStruct {
			structure := arg0.(*Structure)
			return NewBoolean(structure.structType.name == structType.name), nil
		} else {
			return NewBoolean(false), nil
		}
	})
}

func (structType *StructureType) createGetters() []*BuiltinFunc {
	ret := make([]*BuiltinFunc, len(structType.fields))
	for i, field := range structType.fields {
		ret[i] = structType.createGetterForField(field)
	}

	return ret
}

func (structType *StructureType) createGetterForField(field string) *BuiltinFunc {
	getterName := structType.name + "-" + field

	return NewBuiltinFunc(getterName, func(args []ValueObject) (ValueObject, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("getter requires 1 argument, got %v", len(args))
		}

		if args[0].GetValueType() != ValueStruct {
			return nil, fmt.Errorf("getter requires struct as first argument")
		}
		structValue := args[0].(*Structure)

		return structValue.values[field], nil
	})
}

func (structType *StructureType) createSetters() []*BuiltinFunc {
	ret := make([]*BuiltinFunc, len(structType.fields))
	for i, field := range structType.fields {
		ret[i] = structType.createSetterForField(field)
	}

	return ret
}

func (structType *StructureType) createSetterForField(field string) *BuiltinFunc {
	setterName := structType.name + "-set-" + field + "!"

	return NewBuiltinFunc(setterName, func(args []ValueObject) (ValueObject, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("setter requires 2 arguments, got %v", len(args))
		}

		if args[0].GetValueType() != ValueStruct {
			return nil, fmt.Errorf("setter requires struct as first argument")
		}
		structValue := args[0].(*Structure)

		structValue.values[field] = args[1]

		return GetNilObject(), nil
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
