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

func LoadModule(sourceFactory SourceFactory, moduleName, alias string) (map[string]ValueObject, error) {
	parts := strings.Split(moduleName, ModuleSeparator)
	l := len(parts)
	if l == 0 {
		return nil, errors.New(fmt.Sprintf("invalid module %s", moduleName))
	}
	filename := parts[l-1]
	if !strings.HasSuffix(filename, ".boli") {
		filename = filename + ".boli"
	}
	path := strings.Join(append(parts[:l-1], filename), string(os.PathSeparator))

	source, err := sourceFactory.GetSource(path)
	if err != nil {
		return nil, err
	}

	valueMap, err := LoadValues(source)
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

func LoadValues(source Source) (map[string]ValueObject, error) {
	code, err := source.Read()
	if err != nil {
		return nil, err
	}

	interpreter := NewInterpreter(nil)
	_, err = interpreter.Run(code)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]ValueObject)
	for name, entry := range interpreter.env.entries {
		if entry.isOwned {
			ret[name] = entry.value
		}
	}

	return ret, nil
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
