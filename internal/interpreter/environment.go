package interpreter

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

func (env *Environment) Set(name string, value ValueObject) {
	env.values[name] = value
}
