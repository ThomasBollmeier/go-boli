package interpreter

type ValueType int

const (
	ValueNil ValueType = iota
	ValueBoolean
	ValueInteger
	ValueRational
	ValueReal
	ValueString
	ValueBuiltinFunc
	ValueLambda
	ValueInvalid
)

type ValueObject interface {
	GetValueType() ValueType
}
