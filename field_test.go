package graphb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_AddArguments(t *testing.T) {
	f := MakeField("f")
	assert.Equal(t, 0, len(f.Arguments))

	f.AddArguments(ArgumentBool("b", true))
	assert.Equal(t, Argument{"b", argBool(true)}, f.Arguments[0])
}
