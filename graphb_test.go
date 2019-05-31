package graphb

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTheWholePackage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := Query{
			Type: "query",
			Name: "",
			Fields: []*Field{
				{
					Name:  "courses",
					Alias: "Alias",
					Arguments: []Argument{
						ArgumentInt("uid", 123),
						ArgumentStringSlice("blocked_nds", "nd013", "nd014"),
					},
					Fields: Fields("key", "id"),
				},
			},
		}
		str, err := q.String()
		assert.Nil(t, err)
		assert.Equal(t, `query{Alias:courses(uid:123,blocked_nds:["nd013","nd014"]){key,id}}`, str)

		str, err = q.JSON()
		assert.Nil(t, err)
		assert.Equal(t, `{"query":"query{Alias:courses(uid:123,blocked_nds:[\"nd013\",\"nd014\"]){key,id}}"}`, str)

		var buf bytes.Buffer
		err = q.Fields[0].StringChan(&buf)
		assert.Nil(t, err)
		assert.Equal(t, `Alias:courses(uid:123,blocked_nds:["nd013","nd014"]){key,id}`, buf.String())
	})

	t.Run("Invalid names", func(t *testing.T) {
		q := Query{
			Type: "query",
			Name: "test_graphb",
			Fields: []*Field{
				{
					Name:      "courses",
					Alias:     "Lets_Have_An_Alias看",
					Arguments: nil,
					Fields:    Fields("id", "key"),
				},
			},
		}
		strCh, err := q.String()
		assert.Equal(t, "'Lets_Have_An_Alias看' is an invalid alias name in GraphQL. A valid name matches /[_A-Za-z][_0-9A-Za-z]*/, see: http://facebook.github.io/graphql/October2016/#sec-Names", err.Error())
		assert.Equal(t, "", strCh)
	})

	t.Run("check cycles", func(t *testing.T) {
		f := &Field{
			Name:  "courses",
			Alias: "Alias",
			Arguments: []Argument{
				ArgumentInt("uid", 123),
				ArgumentStringSlice("blocked_nds", "nd013", "nd014"),
			},
			Fields: Fields(""),
		}
		f2 := &Field{Fields: Fields("")}
		f2.Fields[0] = f
		f.Fields[0] = f2

		q := Query{
			Type:   "query",
			Name:   "",
			Fields: []*Field{f2},
		}
		strCh, err := q.String()
		assert.IsTypef(t, CyclicFieldErr{}, errors.Cause(err), "")
		assert.Equal(t, "", strCh)

		var buf bytes.Buffer
		err = f.StringChan(&buf)
		assert.IsTypef(t, CyclicFieldErr{}, errors.Cause(err), "")
		assert.Equal(t, "", strCh)
	})

	t.Run("name validation", func(t *testing.T) {
		q := NewQuery(TypeQuery, OfName("我"))
		assert.IsType(t, InvalidNameErr{}, errors.Cause(q.E))

		q = NewQuery(TypeQuery, OfName("_我"))
		assert.IsType(t, InvalidNameErr{}, errors.Cause(q.E))

		q = NewQuery(TypeMutation, OfName("x-x"))
		assert.IsType(t, InvalidNameErr{}, errors.Cause(q.E))

		q = NewQuery(TypeMutation, OfName("x x"))
		assert.IsType(t, InvalidNameErr{}, errors.Cause(q.E))

		q = NewQuery(TypeSubscription, OfName("_1x1_1x1_"))
		assert.Nil(t, q.E)
	})

	t.Run("Nested fields", func(t *testing.T) {
		q := NewQuery(
			TypeQuery,
			OfName("another_test"),
			OfField(
				"users",
				OfFields("id", "username"),
				OfField(
					"threads",
					OfArguments(ArgumentString("title", "A Good Title")),
					OfFields("title", "created_at"),
				),
			),
		)
		assert.Nil(t, q.E)
		s, err := q.JSON()
		assert.Nil(t, err)
		assert.Equal(t, `{"query":"query another_test{users{id,username,threads(title:\"A Good Title\"){title,created_at}}}"}`, s)
	})

	t.Run("Invalid operation type", func(t *testing.T) {
		q := Query{Type: "muTatio"}
		str, err := q.String()
		assert.Equal(t, "'muTatio' is an invalid operation type in GraphQL. A valid type is one of 'query', 'mutation', 'subscription'", err.Error())
		assert.Equal(t, "", str)

		s, err := q.JSON()
		assert.Equal(t, "'muTatio' is an invalid operation type in GraphQL. A valid type is one of 'query', 'mutation', 'subscription'", err.Error())
		assert.Equal(t, "", s)
	})

	t.Run("Nil field error", func(t *testing.T) {
		q := Query{Type: "mutation", Fields: []*Field{nil}}
		ch, err := q.String()
		assert.Equal(t, "nil Field is not allowed. Please initialize a correct Field with NewField(...) function or Field{...} literal", err.Error())
		assert.Equal(t, "", ch)
	})

	t.Run("Nil field error 2", func(t *testing.T) {
		f := Field{Fields: []*Field{nil}}
		err := f.checkCycle()
		assert.IsTypef(t, NilFieldErr{}, errors.Cause(err), "")
	})
}

func TestMethodChaining(t *testing.T) {
	q := MakeQuery(TypeQuery).
		SetName("").
		SetFields(
			MakeField("x").
				SetArguments(
					ArgumentString("string", "123"),
				).
				SetFields(
					MakeField("x").
						SetArguments(
							ArgumentString("string", "123"),
							ArgumentStringSlice("string_slice", "a"),
						),
					MakeField("x"),
				).
				SetAlias("some_alias"),
		).
		AddFields(MakeField("x"))
	s, err := q.JSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"query":"query{some_alias:x(string:\"123\"){x(string:\"123\",string_slice:[\"a\"]),x},x}"}`, s)
}
