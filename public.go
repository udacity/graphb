// graphb is a Graph QL client query builder.

package graphb

import (
	"strings"

	"github.com/pkg/errors"
)

// StringFromChan builds a string from a channel, assuming the channel has been closed.
func StringFromChan(c <-chan string) string {
	var strs []string
	for str := range c {
		strs = append(strs, str)
	}
	return strings.Join(strs, "")
}

///////////////////
// Field Factory //
///////////////////
type FieldOptionInterface interface {
	runFieldOption(f *Field) error
}

// FieldOption implements FieldOptionInterface
type FieldOption func(field *Field) error

func (fco FieldOption) runFieldOption(f *Field) error {
	return fco(f)
}

func NewField(name string, options ...FieldOptionInterface) (*Field, error) {
	f := Field{Name: name}
	for _, op := range options {
		if err := op.runFieldOption(&f); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &f, nil
}

func OfFields(name ...string) FieldOption {
	return func(f *Field) error {
		f.setFields(Fields(name...))
		return nil
	}
}

func OfArguments(arguments ...Argument) FieldOption {
	return func(f *Field) error {
		f.Arguments = arguments
		return nil
	}
}

///////////////////
// Query Factory //
///////////////////
type QueryOptionInterface interface {
	runQueryOption(q *Query) error
}

// QueryOption implements QueryOptionInterface
type QueryOption func(query *Query) error

func (qo QueryOption) runQueryOption(query *Query) error {
	return qo(query)
}

// NewQuery creates a new Query.
// Type and Fields are required.
// Other options such as operation name and alias are optional.
func NewQuery(Type operationType, options ...QueryOptionInterface) (*Query, error) {
	q := &Query{OperationType: Type}

	for _, op := range options {
		if err := op.runQueryOption(q); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return q, nil
}

// OfName returns a QueryOption which validates and sets the operation name of a query.
func OfName(name string) QueryOption {
	return func(query *Query) error {
		if name != "" && !validName.MatchString(name) {
			return errors.Errorf("'%s' is not a valid name.", name)
		}
		query.OperationName = name
		return nil
	}
}

////////////////////////////
// fieldContainer Factory //
////////////////////////////
type fieldContainer interface {
	getFields() []*Field
	setFields([]*Field)
}

// FieldContainerOption implements FieldOptionInterface and QueryOptionInterface.
// FieldContainerOption is a function which takes in a fieldContainer and config it.
// Both Query and Field are fieldContainer,
type FieldContainerOption func(fc fieldContainer) error

func (fco FieldContainerOption) runQueryOption(q *Query) error {
	return fco(q)
}

func (fco FieldContainerOption) runFieldOption(f *Field) error {
	return fco(f)
}

func OfField(name string, options ...FieldOptionInterface) FieldContainerOption {
	return func(fc fieldContainer) error {
		f, err := NewField(name, options...)
		if err != nil {
			return errors.WithStack(err)
		}
		fc.setFields(append(fc.getFields(), f))
		return nil
	}
}

// Fields takes a list of strings and make them a slice of *Field.
// This is useful when you want fields with no sub fields.
// For example:
//	query { courses { id, key } }
// can be written as:
// 	Query{
//		OperationType: "query",
//		Fields: []*Field{
//			{
//				Name:      "courses",
//				Fields:    Fields("id", "key"),
//			},
//		},
//	}
func Fields(args ...string) []*Field {
	fields := make([]*Field, len(args))
	for i, name := range args {
		fields[i] = &Field{
			Name: name,
		}
	}
	return fields
}
