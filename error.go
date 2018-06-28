package graphb

import (
	"fmt"
)

type nameType string

const (
	operationName nameType = "operation name"
	aliasName     nameType = "alias name"
	fieldName     nameType = "field name"
	argumentName  nameType = "argument name"
)

// InvalidNameErr is returned when an invalid name is used. In GraphQL, operation, alias, field and argument all have names.
// A valid name matches ^[_A-Za-z][_0-9A-Za-z]*$ exactly.
type InvalidNameErr struct {
	Type nameType
	Name string
}

func (e InvalidNameErr) Error() string {
	return fmt.Sprintf("'%s' is an invalid %s in GraphQL. A valid name matches /[_A-Za-z][_0-9A-Za-z]*/, see: http://facebook.github.io/graphql/October2016/#sec-Names", e.Name, e.Type)
}

// InvalidOperationTypeErr is returned when the operation is not one of query, mutation and subscription.
type InvalidOperationTypeErr struct {
	Type operationType
}

func (e InvalidOperationTypeErr) Error() string {
	return fmt.Sprintf("'%s' is an invalid operation type in GraphQL. A valid type is one of 'query', 'mutation', 'subscription'", e.Type)
}

// NilFieldErr is returned when any field is nil. Of course the author could choose to ignore nil fields. But, author chose a stricter construct.
type NilFieldErr struct{}

func (e NilFieldErr) Error() string {
	return "nil Field is not allowed. Please initialize a correct Field with NewField(...) function or Field{...} literal"
}

// CyclicFieldErr is returned when any field contains a loop which goes back to itself.
type CyclicFieldErr struct {
	Field Field
}

func (e CyclicFieldErr) Error() string {
	return fmt.Sprintf("Field %+v contains cyclic loop", e.Field)
}

// ArgumentTypeNotSupportedErr is returned when user tries to pass an unsupported type to ArgumentAny.
type ArgumentTypeNotSupportedErr struct {
	Value interface{}
}

func (e ArgumentTypeNotSupportedErr) Error() string {
	return fmt.Sprintf("Argument %+v of Type %T is not supported", e.Value, e.Value)
}
