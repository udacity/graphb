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
