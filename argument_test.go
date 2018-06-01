package graphb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
