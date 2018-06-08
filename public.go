// graphb is a Graph QL client query builder.
// public.go contains public functions (not struct methods) to construct Query(s) and Field(s).

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

// NewField uses functional options to construct a new Field and returns the pointer to it.
// On error, the pointer is nil.
// To know more about this design pattern, see https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewField(name string, options ...FieldOptionInterface) *Field {
	f := &Field{Name: name}
	for _, op := range options {
		if err := op.runFieldOption(f); err != nil {
			f.E = errors.WithStack(err)
			return f
		}
	}
	return f
}

// FieldOptionInterface implements functional options for NewField().
type FieldOptionInterface interface {
	runFieldOption(f *Field) error
}

// FieldOption implements FieldOptionInterface
type FieldOption func(field *Field) error

func (fco FieldOption) runFieldOption(f *Field) error {
	return fco(f)
}

// OfFields returns a FieldOption which sets a list of sub fields of given names of the targeting field.
// All the sub fields only have one level which is their names. That is, no sub fields have sub fields.
func OfFields(name ...string) FieldOption {
	return func(f *Field) error {
		f.setFields(Fields(name...))
		return nil
	}
}

func OfAlias(alias string) FieldOption {
	return func(f *Field) error {
		f.Alias = alias
		return errors.WithStack(f.checkAlias())
	}
}

// OfArguments returns a FieldOption which sets the arguments of the targeting field.
func OfArguments(arguments ...Argument) FieldOption {
	return func(f *Field) error {
		f.Arguments = arguments
		return nil
	}
}

///////////////////
// Query Factory //
///////////////////

// NewQuery uses functional options to construct a new Query and returns the pointer to it.
// On error, the pointer is nil.
// Type is required.
// Other options such as operation name and alias are optional.
func NewQuery(Type operationType, options ...QueryOptionInterface) *Query {
	// todo: change to new style error handling
	q := &Query{Type: Type}

	for _, op := range options {
		if err := op.runQueryOption(q); err != nil {
			q.E = errors.WithStack(err)
			return q
		}
	}
	return q
}

// QueryOptionInterface implements functional options for NewQuery().
type QueryOptionInterface interface {
	runQueryOption(q *Query) error
}

// QueryOption implements QueryOptionInterface
type QueryOption func(query *Query) error

func (qo QueryOption) runQueryOption(query *Query) error {
	return qo(query)
}

// OfName returns a QueryOption which validates and sets the operation name of a query.
func OfName(name string) QueryOption {
	return func(query *Query) error {
		query.Name = name
		if err := query.checkName(); err != nil {
			return errors.WithStack(err)
		}
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

// FieldContainerOption implements FieldOptionInterface and QueryOptionInterface,
// which means, it can be used as the functional option for both NewQuery() and NewField().
// FieldContainerOption is a function which takes in a fieldContainer and config it.
// Both Query and Field are fieldContainer.
type FieldContainerOption func(fc fieldContainer) error

func (fco FieldContainerOption) runQueryOption(q *Query) error {
	return fco(q)
}

func (fco FieldContainerOption) runFieldOption(f *Field) error {
	return fco(f)
}

// OfField returns a FieldContainerOption and has the same parameter signature of
// NewField(name string, options ...FieldOptionInterface) (*Field, error)
func OfField(name string, options ...FieldOptionInterface) FieldContainerOption {
	return func(fc fieldContainer) error {
		f := NewField(name, options...)
		if f.E != nil {
			return errors.WithStack(f.E)
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
//		Type: "query",
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
