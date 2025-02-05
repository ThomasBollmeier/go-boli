package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

func stringConv(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("string: expected 1 argument, got %d", len(objects))
	}

	return NewStr(objects[0].String()), nil
}

func stringSub(objects []ValueObject) (ValueObject, error) {
	numArgs := len(objects)

	if numArgs != 2 && numArgs != 3 {
		return nil, fmt.Errorf("string-sub expects 2 or 3 args, got %d", numArgs)
	}

	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-sub expects string as first argument")
	}
	s := objects[0].(*Str).Value

	if objects[1].GetValueType() != ValueInteger {
		return nil, fmt.Errorf("string-sub expects integer as second argument")
	}
	idx := objects[1].(*Integer).Value
	if idx < 0 || idx >= len(s) {
		return nil, fmt.Errorf("invalid string index")
	}

	var length int
	if numArgs == 3 {
		if objects[2].GetValueType() != ValueInteger {
			return nil, fmt.Errorf("string-sub expects integer as third argument")
		}
		length = objects[2].(*Integer).Value
		if length < 0 {
			return nil, fmt.Errorf("string-sub expects non-negative integer length")
		}
		if length > len(s)-idx {
			length = len(s) - idx
		}

	} else {
		length = len(s) - idx
	}

	return NewStr(s[idx : idx+length]), nil
}

func stringReplace(objects []ValueObject) (ValueObject, error) {
	numArgs := len(objects)

	if numArgs != 3 {
		return nil, fmt.Errorf("string-replace expects 3 args, got %d", numArgs)
	}

	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-replace expects string as first argument")
	}
	s := objects[0].(*Str).Value

	if objects[1].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-replace expects string as second argument")
	}
	oldStr := objects[1].(*Str).Value

	if objects[2].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-replace expects string as third argument")
	}
	newStr := objects[2].(*Str).Value

	return NewStr(strings.Replace(s, oldStr, newStr, -1)), nil
}

func stringConcat(objects []ValueObject) (ValueObject, error) {
	concatenated := ""

	for _, object := range objects {
		if object.GetValueType() != ValueString {
			return nil, fmt.Errorf("string-concat expects strings as arguments")
		}
		concatenated += object.(*Str).Value
	}

	return NewStr(concatenated), nil
}

func stringUpper(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("string-upper expects 1 string argument, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-upper expects string as argument")
	}
	s := objects[0].(*Str).Value

	return NewStr(strings.ToUpper(s)), nil
}

func stringLower(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("string-lower expects 1 string argument, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-lower expects string as argument")
	}
	s := objects[0].(*Str).Value

	return NewStr(strings.ToLower(s)), nil
}

func stringCount(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("string-count expects 1 string argument, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string-count expects string as argument")
	}
	s := objects[0].(*Str).Value

	return NewInteger(len(s)), nil
}

func stringEqual(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 2 {
		return nil, fmt.Errorf("string=? expects 2 string argument, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string=? expects string as first argument")
	}
	s1 := objects[0].(*Str).Value

	if objects[1].GetValueType() != ValueString {
		return nil, fmt.Errorf("string=? expects string as second argument")
	}
	s2 := objects[1].(*Str).Value

	return NewBoolean(s1 == s2), nil
}

func stringToInt(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("string->int expects 1 string argument, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string->int expects string as argument")
	}
	s := objects[0].(*Str).Value

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}

	return NewInteger(i), nil
}

func stringToReal(objects []ValueObject) (ValueObject, error) {
	if len(objects) != 1 {
		return nil, fmt.Errorf("string->int expects 1 string argument, got %d", len(objects))
	}
	if objects[0].GetValueType() != ValueString {
		return nil, fmt.Errorf("string->int expects string as argument")
	}
	s := objects[0].(*Str).Value
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ",", ".", -1)

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}

	return NewReal(f), nil
}
