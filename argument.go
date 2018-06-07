package graphb

import (
	"fmt"
	"strings"
)

type Argument struct {
	Name  string
	Value string
}

func ArgumentAny(name string, value interface{}) (Argument, error) {
	switch v := value.(type) {
	case bool:
		return ArgumentBool(name, v), nil
	case []bool:
		return argumentSlice(name, boolSlice(v)), nil

	case int:
		return ArgumentInt(name, v), nil
	case []int:
		return argumentSlice(name, intSlice(v)), nil

	case string:
		return ArgumentString(name, v), nil
	case []string:
		return argumentSlice(name, stringSlice(v)), nil

	default:
		return Argument{}, ArgumentTypeNotSupportErr{Value: value}
	}
}

func ArgumentBool(name string, value bool) Argument {
	return Argument{name, argBool(value).Rep()}
}

func ArgumentInt(name string, value int) Argument {
	return Argument{name, argInt(value).Rep()}
}

func ArgumentString(name string, value string) Argument {
	return Argument{name, argString(value).Rep()}
}

func ArgumentStringSlice(name string, values ...string) Argument {
	return argumentSlice(name, stringSlice(values))
}

func ArgumentIntSlice(name string, values ...int) Argument {
	return argumentSlice(name, intSlice(values))
}

func ArgumentBoolSlice(name string, values ...bool) Argument {
	return argumentSlice(name, boolSlice(values))
}

//////////////////////////////
// Primitive List Interface //
//////////////////////////////

type valueSlice interface {
	// Get the string representation of i-th element
	Len() int
	String(int) string
}

func argumentSlice(name string, slice valueSlice) Argument {
	representations := make([]string, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		representations[i] = slice.String(i)
	}
	return Argument{Name: name, Value: fmt.Sprintf("[%s]", strings.Join(representations, ","))}
}

//////////////////////////////////
// Primitive List Wrapper Types //
//////////////////////////////////

// boolSlice implements valueSlice
type boolSlice []bool

func (bs boolSlice) Len() int {
	return len(bs)
}

func (bs boolSlice) String(i int) string {
	return argBool(bs[i]).Rep()
}

// stringSlice implements valueSlice
type stringSlice []string

func (s stringSlice) Len() int {
	return len(s)
}

func (s stringSlice) String(i int) string {
	return argString(s[i]).Rep()
}

// intSlice implements valueSlice
type intSlice []int

func (s intSlice) Len() int {
	return len(s)
}

func (s intSlice) String(i int) string {
	return fmt.Sprintf("%v", argInt(s[i]).Rep())
}

/////////////////////////////
// Primitive Wrapper Types //
/////////////////////////////

// argBool represents a boolean value.
type argBool bool

func (v argBool) Rep() string {
	return fmt.Sprintf("%v", v)
}

// argInt represents an integer value.
type argInt int

func (v argInt) Rep() string {
	return fmt.Sprintf("%d", v)
}

// argString represents a string value.
type argString string

func (v argString) Rep() string {
	return fmt.Sprintf(`"%s"`, v)
}
