package interpreter

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const ModuleSeparator string = "::"
const ModulePathEnvVar string = "BOLI_MODULE_PATH"

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
	for _, moduleDir := range getModuleDirs() {
		fullPath := moduleDir + f.path
		content, err := os.ReadFile(fullPath)
		if err == nil {
			return string(content), nil
		}
	}

	return "", errors.New(fmt.Sprintf("File %s not found", f.path))
}

func getModuleDirs() []string {
	moduleDirs := []string{""}

	pathVar, ok := os.LookupEnv(ModulePathEnvVar)
	if ok {
		paths := strings.Split(pathVar, string(os.PathListSeparator))
		for _, path := range paths {
			moduleDirs = append(moduleDirs, path+string(os.PathSeparator))
		}
	}

	return moduleDirs
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

func RunSource(source Source, args []string) (ValueObject, error) {
	code, err := source.Read()
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter(nil)
	value, err := interpreter.Run(code)
	if err != nil {
		return nil, err
	}

	mainValue, ok := interpreter.env.Get("main")
	if !ok {
		return value, nil
	}
	callable, ok := mainValue.(Callable)
	if !ok {
		return nil, errors.New("main is not a callable")
	}
	argValues := make([]ValueObject, 0)
	for _, arg := range args {
		argValues = append(argValues, NewStr(arg))
	}

	return Call(callable, argValues)
}

func makeRequireFn(env *Environment) func(objects []ValueObject) (ValueObject, error) {

	return func(objects []ValueObject) (ValueObject, error) {
		if len(objects) == 0 {
			return nil, fmt.Errorf("require function expects at least one argument")
		}

		var moduleName, alias string
		var valueMap map[string]ValueObject
		var err error

		for _, object := range objects {
			switch object.GetValueType() {
			case ValueSymbol:
				moduleName = object.(*Symbol).Value[1:]
				alias = ""
			case ValuePair:
				pair := object.(*Pair)
				if pair.first.GetValueType() != ValueSymbol || pair.second.GetValueType() != ValueSymbol {
					return nil, fmt.Errorf("require function expects pair of symbols")
				}
				moduleName = pair.first.(*Symbol).Value[1:]
				alias = pair.second.(*Symbol).Value[1:]
			default:
				return nil, fmt.Errorf("require function expects symbols or pair of symbols as argument")
			}

			valueMap, err = loadModule(env.sourceFactory, moduleName, alias)
			if err != nil {
				return nil, err
			}

			for name, value := range valueMap {
				env.Set(name, value, false)
			}
		}

		return GetNilObject(), nil
	}
}

func makeProvideFn(env *Environment) func(objects []ValueObject) (ValueObject, error) {

	return func(objects []ValueObject) (ValueObject, error) {

		if len(objects) == 0 {
			return nil, fmt.Errorf("provide function expects at least one argument")
		}

		var objName, alias string

		for _, object := range objects {
			switch object.GetValueType() {
			case ValueSymbol:
				objName = object.(*Symbol).Value[1:]
				alias = ""
			case ValuePair:
				pair := object.(*Pair)
				if pair.first.GetValueType() != ValueSymbol || pair.second.GetValueType() != ValueSymbol {
					return nil, fmt.Errorf("provide function expects pair of symbols")
				}
				objName = pair.first.(*Symbol).Value[1:]
				alias = pair.second.(*Symbol).Value[1:]
			default:
				return nil, fmt.Errorf("provide function expects symbols or pair of symbols as argument")
			}

			if alias == "" {
				env.providedValues[objName] = objName
			} else {
				env.providedValues[objName] = alias
			}
		}

		return GetNilObject(), nil
	}
}
