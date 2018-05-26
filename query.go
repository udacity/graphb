package graphb

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Query represents a GraphQL query.
// Though all fields (Go struct field, not GraphQL field) of this struct is public,
// the author recommends you to use functions in public.go.
type Query struct {
	OperationType operationType // The operation type is either query, mutation, or subscription.
	OperationName string        // The operation name is a meaningful and explicit name for your operation.
	Fields        []*Field
}

// implements fieldContainer
func (q *Query) getFields() []*Field {
	return q.Fields
}

func (q *Query) setFields(fs []*Field) {
	q.Fields = fs
}

// StringChan returns a string channel and an error.
// When error is not nil, the channel is nil.
// When error is nil, the channel is guaranteed to be closed.
// Warning: One should never receive from a nil channel for eternity awaits by a nil channel.
func (q *Query) StringChan() (<-chan string, error) {
	ch := make(chan string)

	if err := q.check(); err != nil {
		close(ch)
		return ch, errors.WithStack(err)
	}

	for _, f := range q.Fields {
		if f == nil {
			close(ch)
			return ch, errors.WithStack(NilFieldErr{})
		}
		if err := f.check(); err != nil {
			close(ch)
			return ch, errors.WithStack(err)
		}
	}
	return q.stringChan(), nil
}

// StringChan returns a read only channel which is guaranteed to be closed in the future.
func (q *Query) stringChan() <-chan string {
	tokenChan := make(chan string)
	go func() {
		tokenChan <- strings.ToLower(string(q.OperationType))
		// emit operation name
		if q.OperationName != "" {
			tokenChan <- " " // todo: make it a token const, not a raw string
			tokenChan <- q.OperationName
		}
		// emit fields
		tokenChan <- "{"
		for _, field := range q.Fields {
			strs := field.stringChan()
			for str := range strs {
				tokenChan <- str
			}
			tokenChan <- ","
		}
		tokenChan <- "}"
		close(tokenChan)
	}()
	return tokenChan
}

func (q *Query) JsonBody() (string, error) {
	strCh, err := q.StringChan()
	if err != nil {
		return "", errors.WithStack(err)
	}
	s := StringFromChan(strCh)
	return fmt.Sprintf(`{"query":"%s"}`, strings.Replace(s, `"`, `\"`, -1)), nil
}

func (q *Query) check() error {
	// check query
	if !isValidOperationType(q.OperationType) {
		return errors.WithStack(InvalidOperationTypeErr{q.OperationType})
	}
	if q.OperationName != "" && !validName.MatchString(q.OperationName) {
		return errors.WithStack(InvalidNameErr{q.OperationName})
	}
	return nil
}

func (q *Query) SetOperationName(name string) *Query{
	q.OperationName = name
	return q
}

func (q *Query) SetFields(fields ...*Field) *Query {
	q.Fields = fields
	return q
}

func MakeQuery(t operationType) *Query {
	return &Query{OperationType: t}
}