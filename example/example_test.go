package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/udacity/graphb"
)

func TestMethodChaining(t *testing.T) {
	q := graphb.MakeQuery(graphb.TypeQuery).
		SetName("").
		SetFields(
			graphb.MakeField("a").
				SetArguments(
					graphb.ArgumentString("string", "123"),
				).
				SetFields(
					graphb.MakeField("x").
						SetArguments(
							graphb.ArgumentString("string", "123"),
							graphb.ArgumentStringSlice("string_slice", "a"),
						),
					graphb.MakeField("y"),
				).
				SetAlias("some_alias"),
		).
		AddFields(graphb.MakeField("b"))
	s, err := q.JSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"query":"query{some_alias:a(string:\"123\"){x(string:\"123\",string_slice:[\"a\"]),y,},b,}"}`, s)
}
