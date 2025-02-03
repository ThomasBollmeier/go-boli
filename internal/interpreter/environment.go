package interpreter

import (
	"errors"
	"fmt"
)

type Environment struct {
	parent         *Environment
	entries        map[string]EnvEntry
	sourceFactory  SourceFactory
	providedValues map[string]string
}

type EnvEntry struct {
	value   ValueObject
	isOwned bool
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:         parent,
		entries:        make(map[string]EnvEntry),
		sourceFactory:  NewFileSourceFactory(),
		providedValues: make(map[string]string),
	}
}

func NewGlobalEnv() *Environment {
	ret := NewEnvironment(nil)
	for _, op := range []string{"+", "-", "*", "/", "%"} {
		ret.SetGlobalBuiltinFunc(op, makeOperatorFn(op, true))
	}
	ret.SetGlobalBuiltinFunc("^", makeOperatorFn("^", false))
	for _, op := range []string{"=", ">", ">=", "<", "<="} {
		ret.SetGlobalBuiltinFunc(op, func(objects []ValueObject) (ValueObject, error) {
			return compareNumbers(op, objects)
		})
	}
	ret.SetGlobalBuiltinFunc("not", func(objects []ValueObject) (ValueObject, error) {
		if len(objects) != 1 {
			return nil, errors.New("not requires exactly one argument")
		}
		return NewBoolean(!isTruthy(objects[0])), nil
	})

	ret.SetGlobalBuiltinFunc("car", car)
	ret.SetGlobalBuiltinFunc("cdr", cdr)
	ret.SetGlobalBuiltinFunc("pair?", isPair)
	ret.SetGlobalBuiltinFunc("cons", cons)

	ret.SetGlobalBuiltinFunc("list", list)
	ret.SetGlobalBuiltinFunc("list?", isList)
	ret.SetGlobalBuiltinFunc("list-ref", listGetRef)

	ret.SetGlobalBuiltinFunc("vector", vector)
	ret.SetGlobalBuiltinFunc("vector?", isVec)
	ret.SetGlobalBuiltinFunc("vector-ref", vectorGetRef)

	ret.SetGlobalBuiltinFunc("list->stream", listToStream)
	ret.SetGlobalBuiltinFunc("stream?", isStream)

	ret.SetGlobalBuiltinFunc("filter", filter)
	ret.SetGlobalBuiltinFunc("map", mapFunc)
	ret.SetGlobalBuiltinFunc("drop", drop)
	ret.SetGlobalBuiltinFunc("drop-while", dropWhile)
	ret.SetGlobalBuiltinFunc("take", take)

	ret.SetGlobalBuiltinFunc("display", func(objects []ValueObject) (ValueObject, error) {
		if len(objects) != 1 {
			return nil, fmt.Errorf("expected single arg, got %d", len(objects))
		}
		fmt.Print(objects[0])
		return GetNilObject(), nil
	})
	ret.SetGlobalBuiltinFunc("displayln", func(objects []ValueObject) (ValueObject, error) {
		if len(objects) != 1 {
			return nil, fmt.Errorf("expected single arg, got %d", len(objects))
		}
		fmt.Println(objects[0])
		return GetNilObject(), nil
	})

	// module handling:
	ret.SetGlobalBuiltinFunc("require", makeRequireFn(ret))
	ret.SetGlobalBuiltinFunc("provide", makeProvideFn(ret))

	return ret
}

func (env *Environment) Get(name string) (ValueObject, bool) {
	entry, ok := env.entries[name]
	if ok {
		return entry.value, ok
	}
	if env.parent == nil {
		return nil, false
	}
	return env.parent.Get(name)
}

func (env *Environment) GetDefiningEnv(name string) *Environment {
	_, ok := env.entries[name]
	if ok {
		return env
	}
	if env.parent == nil {
		return nil
	}
	return env.parent.GetDefiningEnv(name)
}

func (env *Environment) Set(name string, value ValueObject, isOwned bool) {
	if value.GetValueType() != ValueLambda {
		env.entries[name] = EnvEntry{
			value:   value,
			isOwned: isOwned,
		}
	} else {
		var entry EnvEntry
		var ok bool
		entry, ok = env.entries[name]
		if !ok {
			env.entries[name] = EnvEntry{
				value:   value,
				isOwned: isOwned,
			}
		} else if entry.isOwned == isOwned {
			var existingLambda *LambdaFunc
			existingLambda, ok = entry.value.(*LambdaFunc)
			if ok {
				newLambda := value.(*LambdaFunc)
				_ = existingLambda.Merge(newLambda)
			}
		}
	}
}

func (env *Environment) SetGlobalBuiltinFunc(name string, function func([]ValueObject) (ValueObject, error)) {
	env.Set(name, NewBuiltinFunc(name, function), false)
}

func (env *Environment) GetProvidedValues() (map[string]ValueObject, error) {
	ret := make(map[string]ValueObject)

	if len(env.providedValues) == 0 {
		for name, entry := range env.entries {
			if entry.isOwned {
				ret[name] = entry.value
			}
		}
	} else {
		for name, providedName := range env.providedValues {
			value, ok := env.Get(name)
			if ok {
				ret[providedName] = value
			} else {
				return nil, fmt.Errorf("no value for %s", providedName)
			}
		}
	}

	return ret, nil
}
