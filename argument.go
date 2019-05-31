package graphb

import (
	"bytes"
	"fmt"
)

type argumentValue interface {
	stringChan(buffer *bytes.Buffer)
}

type Argument struct {
	Name  string
	Value argumentValue
}

func (a *Argument) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString(a.Name)
	buffer.WriteString(":")
	a.Value.stringChan(buffer)
}

func ArgumentAny(name string, value interface{}) (Argument, error) {
	switch v := value.(type) {
	case bool:
		return ArgumentBool(name, v), nil
	case []bool:
		return ArgumentBoolSlice(name, v...), nil

	case int:
		return ArgumentInt(name, v), nil
	case []int:
		return ArgumentIntSlice(name, v...), nil

	case string:
		return ArgumentString(name, v), nil
	case []string:
		return ArgumentStringSlice(name, v...), nil

	default:
		return Argument{}, ArgumentTypeNotSupportedErr{Value: value}
	}
}

func ArgumentBool(name string, value bool) Argument {
	return Argument{name, argBool(value)}
}

func ArgumentInt(name string, value int) Argument {
	return Argument{name, argInt(value)}
}

func ArgumentString(name string, value string) Argument {
	return Argument{name, argString(value)}
}

func ArgumentBoolSlice(name string, values ...bool) Argument {
	return Argument{name, argBoolSlice(values)}
}

func ArgumentIntSlice(name string, values ...int) Argument {
	return Argument{name, argIntSlice(values)}
}

func ArgumentStringSlice(name string, values ...string) Argument {
	return Argument{name, argStringSlice(values)}
}

// ArgumentCustomType returns a custom GraphQL type's argument representation, which could be a recursive structure.
func ArgumentCustomType(name string, values ...Argument) Argument {
	return Argument{name, argumentSlice(values)}
}

func ArgumentCustomTypeSlice(name string, values ...[]Argument) Argument {
	return Argument{name, argCustomTypeSlice(values)}
}

func ArgumentCustomTypeSliceElem(values ...Argument) []Argument {
	return values
}

// ///////////////////////////
// Primitive Wrapper Types //
// ///////////////////////////

// argBool represents a boolean value.
type argBool bool

func (v argBool) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf("%t", v))
}

// argInt represents an integer value.
type argInt int

func (v argInt) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf("%d", v))
}

// argString represents a string value.
type argString string

func (v argString) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`"%s"`, v))
}

// ////////////////////////////////
// Primitive List Wrapper Types //
// ////////////////////////////////

// argBoolSlice implements valueSlice
type argBoolSlice []bool

func (s argBoolSlice) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString("[")
	for i, v := range s {
		if i != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(fmt.Sprintf("%t", v))
	}
	buffer.WriteString("]")
}

// argIntSlice implements valueSlice
type argIntSlice []int

func (s argIntSlice) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString("[")
	for i, v := range s {
		if i != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(fmt.Sprintf("%d", v))
	}
	buffer.WriteString("]")
}

// argStringSlice implements valueSlice
type argStringSlice []string

func (s argStringSlice) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString("[")
	for i, v := range s {
		if i != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(fmt.Sprintf(`"%s"`, v))
	}
	buffer.WriteString("]")
}

type argumentSlice []Argument

func (s argumentSlice) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString("{")
	for i, v := range s {
		if i != 0 {
			buffer.WriteString(",")
		}
		v.stringChan(buffer)
	}
	buffer.WriteString("}")
}

type argCustomTypeSlice [][]Argument

func (s argCustomTypeSlice) stringChan(buffer *bytes.Buffer) {
	buffer.WriteString("[")
	for i, v := range s {
		if i != 0 {
			buffer.WriteString(",")
		}
		argumentSlice(v).stringChan(buffer)
	}
	buffer.WriteString("]")
}
