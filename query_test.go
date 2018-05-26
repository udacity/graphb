package graphb

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestQuery_checkName(t *testing.T) {
	q := Query{Name: "1"}
	err := q.checkName()
	assert.IsType(t, InvalidNameErr{}, errors.Cause(err))
	assert.Equal(t, "'1' is an invalid operation name in GraphQL. A valid name matches /[_A-Za-z][_0-9A-Za-z]*/, see: http://facebook.github.io/graphql/October2016/#sec-Names", err.Error())
}

func TestQuery_check(t *testing.T) {
	q := Query{Name: "1", Type: TypeQuery}
	err := q.check()
	assert.IsType(t, InvalidNameErr{}, errors.Cause(err))
	assert.Equal(t, "'1' is an invalid operation name in GraphQL. A valid name matches /[_A-Za-z][_0-9A-Za-z]*/, see: http://facebook.github.io/graphql/October2016/#sec-Names", err.Error())
}
