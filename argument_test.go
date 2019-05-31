package graphb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentAny(t *testing.T) {
	arg, err := ArgumentAny("arg", 1)
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", argInt(1)}, arg)

	arg, err = ArgumentAny("arg", true)
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", argBool(true)}, arg)

	arg, err = ArgumentAny("arg", "str")
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", argString("str")}, arg)

	arg, err = ArgumentAny("arg", []string{"str", "slice"})
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", argStringSlice([]string{"str", "slice"})}, arg)

	arg, err = ArgumentAny("arg", []bool{true, false})
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", argBoolSlice([]bool{true, false})}, arg)

	arg, err = ArgumentAny("arg", []int{1, 2})
	assert.Nil(t, err)
	assert.Equal(t, Argument{"arg", argIntSlice([]int{1, 2})}, arg)

	// Type Not Supported
	arg, err = ArgumentAny("arg", 1.1)
	assert.IsType(t, ArgumentTypeNotSupportedErr{}, err)
	assert.Equal(t, "Argument 1.1 of Type float64 is not supported", err.Error())
	assert.Equal(t, Argument{}, arg)
}

func TestArgumentBool(t *testing.T) {
	a := ArgumentBool("blocked", true)
	assert.Equal(t, Argument{"blocked", argBool(true)}, a)

	a = ArgumentBool("blocked", false)
	assert.Equal(t, Argument{"blocked", argBool(false)}, a)
}

func TestArgumentInt(t *testing.T) {
	a := ArgumentInt("blocked", 1)
	assert.Equal(t, Argument{"blocked", argInt(1)}, a)
}

func TestArgumentString(t *testing.T) {
	a := ArgumentString("blocked", "a")
	assert.Equal(t, Argument{"blocked", argString("a")}, a)
	a = ArgumentString("blocked", "")
	assert.Equal(t, Argument{"blocked", argString("")}, a)
}

func TestArgumentStringSlice(t *testing.T) {
	a := ArgumentStringSlice("blocked", "a", "b", "", " ", "d")
	assert.Equal(t, Argument{"blocked", argStringSlice([]string{"a", "b", "", " ", "d"})}, a)

	a = ArgumentStringSlice("blocked")
	assert.Equal(t, Argument{"blocked", argStringSlice(nil)}, a)
}

func Test_argBool(t *testing.T) {
	b := argBool(true)
	i := 0
	var buf bytes.Buffer
	b.stringChan(&buf)
	assert.Equal(t, "true", buf.String())
	i++
	assert.Equal(t, 1, i)
}

func Test_argBoolSlice(t *testing.T) {
	bs := argBoolSlice([]bool{true, false})
	var buf bytes.Buffer
	bs.stringChan(&buf)
	tokens := "[true,false]"
	assert.Equal(t, tokens, buf.String())
}

func Test_argIntSlice(t *testing.T) {
	is := argIntSlice([]int{1, -1, 0})
	tokens := "[1,-1,0]"
	var buf bytes.Buffer
	is.stringChan(&buf)
	assert.Equal(t, tokens, buf.String())
}
