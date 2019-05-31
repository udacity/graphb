package graphb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Query represents a GraphQL query.
// Though all fields (Go struct field, not GraphQL field) of this struct is public,
// the author recommends you to use functions in public.go.
type Query struct {
	Type   operationType // The operation type is either query, mutation, or subscription.
	Name   string        // The operation name is a meaningful and explicit name for your operation.
	Fields []*Field
	E      error
}

// implements fieldContainer
func (q *Query) getFields() []*Field {
	return q.Fields
}

func (q *Query) setFields(fs []*Field) {
	q.Fields = fs
}

// String returns a string channel and an error.
// When error is not nil, the channel is nil.
// When error is nil, the channel is guaranteed to be closed.
// Warning: One should never receive from a nil channel for eternity awaits by a nil channel.
func (q *Query) String() (string, error) {
	var buffer bytes.Buffer
	if err := q.check(); err != nil {
		return "", errors.WithStack(err)
	}

	for _, f := range q.Fields {
		if f == nil {
			return "", errors.WithStack(NilFieldErr{})
		}
		if err := f.check(); err != nil {
			return "", errors.WithStack(err)
		}
	}
	q.string(&buffer)
	return buffer.String(), nil
}

// String returns a read only channel which is guaranteed to be closed in the future.
func (q *Query) string(buffer *bytes.Buffer) {
	buffer.WriteString(strings.ToLower(string(q.Type)))
	// emit operation name
	if q.Name != "" {
		buffer.WriteString(tokenSpace)
		buffer.WriteString(q.Name)
	}
	// emit fields
	buffer.WriteString(tokenLB)
	for i, field := range q.Fields {
		if i != 0 {
			buffer.WriteString(tokenComma)
		}
		field.stringChan(buffer)
	}
	buffer.WriteString(tokenRB)
}

func (q *Query) check() error {
	// check query
	if !isValidOperationType(q.Type) {
		return errors.WithStack(InvalidOperationTypeErr{q.Type})
	}
	if err := q.checkName(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (q *Query) checkName() error {
	if q.Name != "" && !validName.MatchString(q.Name) {
		return errors.WithStack(InvalidNameErr{operationName, q.Name})
	}
	return nil
}

// //////////////
// Public API //
// //////////////

// MakeQuery constructs a Query of the given type and returns a pointer of it.
func MakeQuery(Type operationType) *Query {
	return &Query{Type: Type}
}

// JSON returns a json string with "query" field.
func (q *Query) JSON() (string, error) {
	strCh, err := q.String()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return fmt.Sprintf(`{"query":"%s"}`, strings.Replace(strCh, `"`, `\"`, -1)), nil
}

// SetName sets the Name field of this Query.
func (q *Query) SetName(name string) *Query {
	q.Name = name
	return q
}

// GetField return the field identified by the name. Nil if not exist.
func (q *Query) GetField(name string) *Field {
	for _, f := range q.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// SetFields sets the Fields field of this Query.
// If q.Fields already contains data, they will be replaced.
func (q *Query) SetFields(fields ...*Field) *Query {
	q.Fields = fields
	return q
}

// AddFields adds to the Fields field of this Query.
func (q *Query) AddFields(fields ...*Field) *Query {
	q.Fields = append(q.Fields, fields...)
	return q
}
