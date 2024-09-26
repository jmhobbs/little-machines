package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeComments(t *testing.T) {
	t.Run("no comments", func(t *testing.T) {
		s := "this is a test"
		assert.Equal(t, s, removeComments(s))
	})

	t.Run("only comments", func(t *testing.T) {
		s := "; this is a comment"
		assert.Equal(t, "", removeComments(s))
	})

	t.Run("comments at end", func(t *testing.T) {
		s := "this is a test ; this is a comment"
		assert.Equal(t, "this is a test", removeComments(s))
	})
}

func Test_parseLine(t *testing.T) {
	t.Run("no operands", func(t *testing.T) {
		s := "CLS"
		i, o := parseLine(s)
		assert.Equal(t, "CLS", i)
		assert.Empty(t, o)
	})

	t.Run("one operand", func(t *testing.T) {
		s := "JMP 1020"
		i, o := parseLine(s)
		assert.Equal(t, "JMP", i)
		assert.Equal(t, []string{"1020"}, o)
	})

	t.Run("two operands", func(t *testing.T) {
		s := "LD V1, V2"
		i, o := parseLine(s)
		assert.Equal(t, "LD", i)
		assert.Equal(t, []string{"V1", "V2"}, o)
	})

	t.Run("casts to uppercase", func(t *testing.T) {
		s := "ld v1, v2"
		i, o := parseLine(s)
		assert.Equal(t, "LD", i)
		assert.Equal(t, []string{"V1", "V2"}, o)
	})
}
