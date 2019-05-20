package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/udacity/graphb"
)

func TestMethodChaining(t *testing.T) {
	/*
		This is a method chaining example.

		MakeQuery() creates the query object.

		All the SetXXX methods set the corresponding properties and return a pointer to the original query.

		The order of method calling does not matter. Of course the latter will override the former if the same method.

		Method chaining has the nice feature of IDE code suggestion, so that you don't have to remember all methods.
	*/
	q := graphb.MakeQuery(graphb.TypeQuery).
		SetName("").
		SetFields(
			graphb.MakeField("a").
				SetArguments(
					graphb.ArgumentString("string", "123"),
					graphb.ArgumentCustomTypeSlice(
						"mapArray",
						graphb.ArgumentCustomTypeSliceElem(
							graphb.ArgumentString("foo", "bar"),
							graphb.ArgumentInt("fizzbuzz", 15),
						),
						graphb.ArgumentCustomTypeSliceElem(
							graphb.ArgumentString("foo", "baz"),
							graphb.ArgumentInt("fizzbuzz", 45),
						),
					),
				).
				SetFields(
					graphb.MakeField("x").
						SetArguments(
							graphb.ArgumentString("string", "123"),
							graphb.ArgumentIntSlice("int_slice", 1, 2, 3),
						),
					graphb.MakeField("y"),
				).
				SetAlias("some_alias"),
		).
		AddFields(graphb.MakeField("b"))
	s, err := q.JSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"query":"query{some_alias:a(string:\"123\",mapArray:[{foo:\"bar\",fizzbuzz:15},{foo:\"baz\",fizzbuzz:45}]){x(string:\"123\",int_slice:[1,2,3]),y},b}"}`, s)
}

func TestFunctionalOptions(t *testing.T) {
	/*
		This is a functional option example.

		Depending on how you view it, it may look less verbose than method chaining.

		But, it doesn't have the nice code suggestions.

		I use the naming convention OfSomething() for all functional option function.
	*/
	q := graphb.NewQuery(
		graphb.TypeQuery,
		graphb.OfName("Good_Name_Is_Important"),
		graphb.OfField(
			"books",
			graphb.OfArguments(
				graphb.ArgumentString("author", "William Shakespeare"),
				graphb.ArgumentStringSlice("title", "Hamlet", "Henry IV"),
			),
			graphb.OfFields("author", "title", "price"),
		),
	)
	jsonString, err := q.JSON()
	assert.Nil(t, err)
	assert.Equal(
		t,
		`{"query":"query Good_Name_Is_Important{books(author:\"William Shakespeare\",title:[\"Hamlet\",\"Henry IV\"]){author,title,price}}"}`,
		jsonString,
	)
}

func TestStructLiteral(t *testing.T) {
	/*
		You might consider struct literal more readable, it is definitely more explicit.
	*/
	q := graphb.Query{
		Type: graphb.TypeMutation,
		Name: "It_is_A_Name",
		Fields: []*graphb.Field{
			{
				Name:      "f_1",
				Alias:     "f_1_alias",
				Arguments: []graphb.Argument{graphb.ArgumentIntSlice("arg", 3, 2, -1)},
			},
			graphb.NewField("f_2"),
		},
	}

	jsonString, err := q.JSON()
	assert.Nil(t, err)
	assert.Equal(
		t,
		`{"query":"mutation It_is_A_Name{f_1_alias:f_1(arg:[3,2,-1]),f_2}"}`,
		jsonString,
	)
}
