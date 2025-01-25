package interpreter

import "fmt"

type Environment struct {
	parent *Environment
	values map[string]ValueObject
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		values: make(map[string]ValueObject),
	}
}

func NewGlobalEnv() *Environment {
	ret := NewEnvironment(nil)
	for _, op := range []string{"+", "-", "*", "/", "%"} {
		ret.Set(op, NewBuiltinFunc(op, makeOperatorFn(op, true)))
	}
	ret.Set("^", NewBuiltinFunc("^", makeOperatorFn("^", false)))
	for _, op := range []string{"=", ">", ">=", "<", "<="} {
		ret.Set(op, NewBuiltinFunc(op, func(objects []ValueObject) (ValueObject, error) {
			return compareNumbers(op, objects)
		}))
	}

	ret.SetBuiltinFunc("car", car)
	ret.SetBuiltinFunc("cdr", cdr)
	ret.SetBuiltinFunc("pair?", isPair)
	ret.SetBuiltinFunc("list?", isList)
	ret.SetBuiltinFunc("cons", cons)
	ret.SetBuiltinFunc("list", list)
	ret.SetBuiltinFunc("list-ref", listGetRef)

	ret.SetBuiltinFunc("displayln", func(objects []ValueObject) (ValueObject, error) {
		if len(objects) != 1 {
			return nil, fmt.Errorf("expected single arg, got %d", len(objects))
		}
		fmt.Println(objects[0])
		return GetNilObject(), nil
	})

	return ret
}

func (env *Environment) Get(name string) (ValueObject, bool) {
	value, ok := env.values[name]
	if ok {
		return value, ok
	}
	if env.parent == nil {
		return nil, false
	}
	return env.parent.Get(name)
}

func (env *Environment) GetDefiningEnv(name string) *Environment {
	_, ok := env.values[name]
	if ok {
		return env
	}
	if env.parent == nil {
		return nil
	}
	return env.parent.GetDefiningEnv(name)
}

func (env *Environment) Set(name string, value ValueObject) {
	env.values[name] = value
}

func (env *Environment) SetBuiltinFunc(name string, function func([]ValueObject) (ValueObject, error)) {
	env.Set(name, NewBuiltinFunc(name, function))
}
