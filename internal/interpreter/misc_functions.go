package interpreter

import (
	"fmt"
)

func readLine(objects []ValueObject) (ValueObject, error) {
	var prompt string
	switch len(objects) {
	case 0:
		prompt = ""
	case 1:
		if objects[0].GetValueType() != ValueString {
			return nil, fmt.Errorf("read-line expects string as first argument")
		}
		prompt = objects[0].(*Str).Value
	default:
		return nil, fmt.Errorf("read-line expects at most one argument, got %d", len(objects))
	}

	var line string
	if prompt != "" {
		fmt.Print(prompt)
	}
	_, err := fmt.Scanln(&line)
	if err != nil {
		return nil, err
	}

	return NewStr(line), nil
}

func display(objects []ValueObject) (ValueObject, error) {
	s := ""
	for _, object := range objects {
		s += object.String()
	}
	fmt.Print(s)

	return GetNilObject(), nil
}

func displayln(objects []ValueObject) (ValueObject, error) {
	s := ""
	for _, object := range objects {
		s += object.String()
	}
	fmt.Println(s)

	return GetNilObject(), nil
}

func write(objects []ValueObject) (ValueObject, error) {
	s := ""
	for _, object := range objects {
		switch object.GetValueType() {
		case ValueString:
			s += object.(*Str).Value
		default:
			s += object.String()
		}
	}
	fmt.Print(s)

	return GetNilObject(), nil
}

func writeln(objects []ValueObject) (ValueObject, error) {
	s := ""
	for _, object := range objects {
		switch object.GetValueType() {
		case ValueString:
			s += object.(*Str).Value
		default:
			s += object.String()
		}
	}
	fmt.Println(s)

	return GetNilObject(), nil
}
