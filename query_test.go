package graphb

import (
	"strings"
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

func TestQuery_GetField(t *testing.T) {
	q := MakeQuery(TypeQuery).SetFields(MakeField("f1"))
	f := q.GetField("f1")
	assert.Equal(t, "f1", f.Name)

	f = q.GetField("f2")
	assert.Nil(t, f)
}

func TestQuery_JSON(t *testing.T) {
	t.Parallel()

	t.Run("Arguments can be nested structures", func(t *testing.T) {
		t.Parallel()

		q := NewQuery(TypeMutation).
			SetFields(
				NewField("createQuestion").
					SetArguments(
						ArgumentCustomType(
							"input",
							ArgumentString("title", "what"),
							ArgumentString("content", "what"),
							ArgumentStringSlice("tagIds"),
						),
					).
					SetFields(
						NewField("question", OfFields("id")),
					),
			)

		c := q.stringChan()

		var strs []string
		for str := range c {
			strs = append(strs, str)
		}

		assert.Equal(t, `mutation{createQuestion(input:{title:"what",content:"what",tagIds:[]}){question{id}}}`, strings.Join(strs, ""))
	})

}
