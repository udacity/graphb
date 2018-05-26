package graphb

import (
	"regexp"
	"strings"
)

// checks the validity of a name according to the spec: http://facebook.github.io/graphql/October2016/#sec-Names
var validName = regexp.MustCompile("^[_A-Za-z][_0-9A-Za-z]*$")

func isValidOperationType(Type operationType) bool {
	low := strings.ToLower(string(Type))
	return low == "query" || low == "mutation" || low == "subscription"
}
