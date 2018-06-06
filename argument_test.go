package graphb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentAny(t *testing.T) {
	arg, err := ArgumentAny("arg", 1)
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", "1"}, arg)

	arg, err = ArgumentAny("arg", true)
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", "true"}, arg)

	arg, err = ArgumentAny("arg", "str")
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", `"str"`}, arg)

	arg, err = ArgumentAny("arg", []string{"str", "slice"})
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", `["str","slice"]`}, arg)

	arg, err = ArgumentAny("arg", []bool{true, false})
	assert.Equalf(t, "Argument [true false] of Type []bool is not supported", err.Error(), "This type is not supported yet")
	assert.Equal(t, Argument{}, arg)
}

func TestArgumentBool(t *testing.T) {
	a := ArgumentBool("blocked", true)
	assert.Equal(t, Argument{"blocked", "true"}, a)

	a = ArgumentBool("blocked", false)
	assert.Equal(t, Argument{"blocked", "false"}, a)
}

func TestArgumentInt(t *testing.T) {
	a := ArgumentInt("blocked", 1)
	assert.Equal(t, Argument{"blocked", "1"}, a)
}

func TestArgumentString(t *testing.T) {
	a := ArgumentString("blocked", "a")
	assert.Equal(t, Argument{"blocked", `"a"`}, a)
	a = ArgumentString("blocked", "")
	assert.Equal(t, Argument{"blocked", `""`}, a)
}

func TestArgumentStringSlice(t *testing.T) {
	a := ArgumentStringSlice("blocked", "a", "b", "", " ", "d")
	assert.Equal(t, Argument{"blocked", `["a","b",""," ","d"]`}, a)

	a = ArgumentStringSlice("blocked")
	assert.Equal(t, Argument{"blocked", "[]"}, a)
}
