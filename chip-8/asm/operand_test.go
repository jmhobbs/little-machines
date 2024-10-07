package asm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_register(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var i uint8
		for i = 0; i < 16; i++ {
			r, err := register(fmt.Sprintf("V%x", i))
			assert.Nil(t, err)
			assert.Equal(t, i, r)
		}
	})

	t.Run("invalid: lowercase", func(t *testing.T) {
		_, err := register("v1")
		assert.NotNil(t, err)
	})

	t.Run("invalid: too short", func(t *testing.T) {
		_, err := register("V")
		assert.NotNil(t, err)
	})

	t.Run("invalid: too high", func(t *testing.T) {
		_, err := register("V10")
		assert.NotNil(t, err)
	})

	t.Run("invalid: not a register", func(t *testing.T) {
		_, err := register("x1")
		assert.NotNil(t, err)
	})
}
