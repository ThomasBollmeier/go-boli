package interpreter

type BuiltInModule struct {
	valuesMap map[string]ValueObject
}

func NewBuiltInModule() *BuiltInModule {
	return &BuiltInModule{
		valuesMap: make(map[string]ValueObject),
	}
}

func (m *BuiltInModule) AddValue(name string, value ValueObject) {
	m.valuesMap[name] = value
}

func (m *BuiltInModule) GetValues() map[string]ValueObject {
	return m.valuesMap
}

type BuiltInModules map[string]*BuiltInModule

var builtinModules = createModules()

func createModules() BuiltInModules {
	result := make(BuiltInModules)
	result["math"] = NewMath()
	return result
}
