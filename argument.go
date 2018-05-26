package graphb

import (
	"fmt"
	"strings"
)

// Argument is a size 2 array.
// The first element is the name, the second element is the value of a GraphQL field argument.
type Argument [2]string

func ArgumentInt(name string, value int) Argument {
	return Argument{name, fmt.Sprintf("%d", value)}
}

func ArgumentString(name string, value string) Argument {
	return Argument{name, fmt.Sprintf(`"%s"`, value)}
}

func ArgumentStringSlice(name string, args ...string) Argument {
	return Argument{name, fmt.Sprintf("[%s]", strings.Join(mapStringSliceToStrRepSlice(args), ","))}
}
