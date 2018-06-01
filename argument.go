package graphb

import (
	"fmt"
	"strings"
)

// Argument is a size 2 array.
// The first element is the name.
// The second element is the value of a GraphQL field argument, which is a string.
type Argument struct {
	Name  string
	Value string
}

func ArgumentAny(name string, value interface{}) (Argument, error) {
	switch v := value.(type) {
	case bool:
		return ArgumentBool(name, v), nil
	case int:
		return ArgumentInt(name, v), nil
	case string:
		return ArgumentString(name, v), nil
	case []string:
		return ArgumentStringSlice(name, v...), nil
	default:
		return Argument{}, ArgumentTypeNotSupportErr{Value: value}
	}
}

func ArgumentBool(name string, value bool) Argument {
	return Argument{name, fmt.Sprintf("%v", value)}
}

func ArgumentInt(name string, value int) Argument {
	return Argument{name, fmt.Sprintf("%d", value)}
}

func ArgumentString(name string, value string) Argument {
	return Argument{name, fmt.Sprintf(`"%s"`, value)}
}

// ArgumentStringSlice returns an argument with string list value.
// If the caller wants an empty list [],
// do:
// ArgumentStringSlice("name")
func ArgumentStringSlice(name string, args ...string) Argument {
	return Argument{name, fmt.Sprintf("[%s]", strings.Join(mapStringSliceToStrRepSlice(args), ","))}
}
