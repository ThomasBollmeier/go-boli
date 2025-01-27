package interpreter

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const ModuleSeparator string = "::"

type Source interface {
	Read() (string, error)
}

type SourceFactory interface {
	GetSource(path string) (Source, error)
}

type FileSource struct {
	path string
}

func NewFileSource(path string) *FileSource {
	return &FileSource{path: path}
}

func (f *FileSource) Read() (string, error) {
	content, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

type FileSourceFactory struct{}

func NewFileSourceFactory() *FileSourceFactory {
	return &FileSourceFactory{}
}

func (f *FileSourceFactory) GetSource(path string) (Source, error) {

	return NewFileSource(path), nil
}

func loadModule(sourceFactory SourceFactory, moduleName, alias string) (map[string]ValueObject, error) {
	parts := strings.Split(moduleName, ModuleSeparator)
	l := len(parts)
	if l == 0 {
		return nil, errors.New(fmt.Sprintf("invalid module %s", moduleName))
	}
	filename := parts[l-1] + ".boli"
	path := strings.Join(append(parts[:l-1], filename), string(os.PathSeparator))

	source, err := sourceFactory.GetSource(path)
	if err != nil {
		return nil, err
	}

	valueMap, err := loadValues(source)
	if err != nil {
		return nil, err
	}

	if alias == "" {
		return valueMap, nil
	}

	aliasedValueMap := make(map[string]ValueObject)
	for name, value := range valueMap {
		aliasedValueMap[alias+ModuleSeparator+name] = value
	}

	return aliasedValueMap, nil
}

func loadValues(source Source) (map[string]ValueObject, error) {
	code, err := source.Read()
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter(nil)
	_, err = interpreter.Run(code)
	if err != nil {
		return nil, err
	}

	providedValues, err := interpreter.env.GetProvidedValues()
	if err != nil {
		return nil, err
	}

	return providedValues, nil
}

func RunSource(source Source) (ValueObject, error) {
	code, err := source.Read()
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter(nil)
	value, err := interpreter.Run(code)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func makeRequireFn(env *Environment) func(objects []ValueObject) (ValueObject, error) {

	return func(objects []ValueObject) (ValueObject, error) {
		l := len(objects)

		if l != 1 && l != 2 {
			return nil, fmt.Errorf("require function expects one or two arguments (got %d)", len(objects))
		}

		if objects[0].GetValueType() != ValueSymbol {
			return nil, fmt.Errorf("require function expects symbol")
		}
		moduleName := objects[0].(*Symbol).Value[1:]

		alias := ""
		if l == 2 {
			if objects[1].GetValueType() != ValueSymbol {
				return nil, fmt.Errorf("require function expects symbol")
			}
			alias = objects[1].(*Symbol).Value[1:]
		}

		valueMap, err := loadModule(env.sourceFactory, moduleName, alias)
		if err != nil {
			return nil, err
		}

		for name, value := range valueMap {
			env.Set(name, value, false)
		}

		return GetNilObject(), nil
	}
}

func makeProvideFn(env *Environment) func(objects []ValueObject) (ValueObject, error) {

	return func(objects []ValueObject) (ValueObject, error) {
		l := len(objects)

		if l != 1 && l != 2 {
			return nil, fmt.Errorf("provide function expects one or two arguments (got %d)", len(objects))
		}

		if objects[0].GetValueType() != ValueSymbol {
			return nil, fmt.Errorf("provide function expects symbol")
		}
		objName := objects[0].(*Symbol).Value[1:]

		alias := ""
		if l == 2 {
			if objects[1].GetValueType() != ValueSymbol {
				return nil, fmt.Errorf("provide function expects symbol")
			}
			alias = objects[1].(*Symbol).Value[1:]
		}

		if alias == "" {
			env.providedValues[objName] = objName
		} else {
			env.providedValues[objName] = alias
		}

		return GetNilObject(), nil
	}
}
