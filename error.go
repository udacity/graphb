package graphb

import "fmt"

type nameType string

const (
	operationName nameType = "operation name"
	aliasName     nameType = "alias name"
	fieldName     nameType = "field name"
	argumentName  nameType = "argument name"
)

type InvalidNameErr struct {
	Type nameType
	Name string
}

func (e InvalidNameErr) Error() string {
	return fmt.Sprintf("'%s' is an invalid %s in GraphQL. A valid name matches /[_A-Za-z][_0-9A-Za-z]*/, see: http://facebook.github.io/graphql/October2016/#sec-Names", e.Name, e.Type)
}

type InvalidOperationTypeErr struct {
	Type operationType
}

func (e InvalidOperationTypeErr) Error() string {
	return fmt.Sprintf("'%s' is an invalid operation type in GraphQL. A valid type is one of 'query', 'mutation', 'subscription'", e.Type)
}

type NilFieldErr struct{}

func (e NilFieldErr) Error() string {
	return "nil Field is not allowed. Please initialize a correct Field with NewField(...) function or Field{...} literal"
}

type CyclicFieldErr struct {
	Field Field
}

func (e CyclicFieldErr) Error() string {
	return fmt.Sprintf("Field %+v contains cyclic loop", e.Field)
}
