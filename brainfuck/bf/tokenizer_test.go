package bf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tokenize(t *testing.T) {
	program := bytes.NewReader([]byte("\n><+-.,[] # comment goes here\n> < + - . , [ ]\n"))

	tokens, err := Tokenize(program)
	assert.NoError(t, err)

	assert.Equal(t, []Token{
		MoveRight, MoveLeft, Increment, Decrement, Output, Input, JumpForward, JumpBackward,
		MoveRight, MoveLeft, Increment, Decrement, Output, Input, JumpForward, JumpBackward,
	}, tokens)
}
