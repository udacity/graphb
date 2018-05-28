package graphb

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestOfAlias(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := OfAlias("a").runFieldOption(&Field{})
		assert.Nil(t, err)
	})
	t.Run("failure", func(t *testing.T) {
		err := OfAlias("123").runFieldOption(&Field{})
		assert.IsType(t, InvalidNameErr{}, errors.Cause(err))
	})
}
