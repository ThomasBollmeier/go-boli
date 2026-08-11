package interpreter

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type MockSourceFactory struct {
	sources map[string]Source
}

func NewMockSourceFactory() *MockSourceFactory {
	return &MockSourceFactory{
		sources: make(map[string]Source),
	}
}

func (f *MockSourceFactory) addSource(path string, source Source) {
	f.sources[path] = source
}

func (f *MockSourceFactory) GetSource(path string) (Source, error) {
	source, ok := f.sources[path]
	if !ok {
		return nil, fmt.Errorf("no such source %s", path)
	}
	return source, nil
}

type MockSource struct {
	code string
}

func NewMockSource(code string) *MockSource {
	return &MockSource{code: code}
}

func (f *MockSource) Read() (string, error) {
	return f.code, nil
}

func TestRequire(t *testing.T) {
	factory := NewMockSourceFactory()

	libCode := `
		(def (my-add xs...)
			(+ ...xs))`
	libSource := NewMockSource(libCode)
	factory.addSource("utils.boli", libSource)

	mainCode := `
		(require 'utils)
		(def (run)
			(my-add 40 2))
		(run)`

	expected := &Integer{42}

	interpreter := NewInterpreter(nil)
	interpreter.env.sourceFactory = factory // <-- inject mock

	actual, err := interpreter.Run(mainCode)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, &Integer{42}) {
		t.Errorf("Run() got = %v, want %v", actual, expected)
	}

}

func TestRequireBuiltin(t *testing.T) {
	mainCode := `
		(require ('math . 'm))
		(m::sqrt 16,0)`

	expected := &Real{4.0}

	interpreter := NewInterpreter(nil)
	interpreter.env.sourceFactory = NewMockSourceFactory()

	actual, err := interpreter.Run(mainCode)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Run() got = %v, want %v", actual, expected)
	}
}

func TestRequireNestedWithAlias(t *testing.T) {
	factory := NewMockSourceFactory()

	libCode := `
		(def (my-add xs...)
			(+ ...xs))`
	libSource := NewMockSource(libCode)
	factory.addSource(strings.Join([]string{"utils", "nums.boli"}, string(os.PathSeparator)), libSource)

	mainCode := `
		(require ('utils::nums . 'n))
		(def (run)
			(n::my-add 40 2))
		(run)`

	expected := &Integer{42}

	interpreter := NewInterpreter(nil)
	interpreter.env.sourceFactory = factory // <-- inject mock

	actual, err := interpreter.Run(mainCode)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, &Integer{42}) {
		t.Errorf("Run() got = %v, want %v", actual, expected)
	}

}

func TestProvide(t *testing.T) {
	factory := NewMockSourceFactory()

	libCode := `
		(provide
			'my-add
			('mk-adder . 'make-adder))
		(def (mk-adder n)
			(λ (m) (my-add n m)))
		(def (my-add xs...)
			(+ ...xs))`
	libSource := NewMockSource(libCode)
	factory.addSource("utils.boli", libSource)

	mainCode := `
		(require 'utils)
		(def (run)
			((make-adder (my-add 1 1)) 40))
		(run)`

	expected := &Integer{42}

	interpreter := NewInterpreter(nil)
	interpreter.env.sourceFactory = factory // <-- inject mock

	actual, err := interpreter.Run(mainCode)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, &Integer{42}) {
		t.Errorf("Run() got = %v, want %v", actual, expected)
	}

}
