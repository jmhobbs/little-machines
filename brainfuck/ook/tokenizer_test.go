package ook

import (
	"bytes"
	"testing"

	"github.com/jmhobbs/little-machines/brainfuck/bf"
	"github.com/stretchr/testify/assert"
)

func Test_Tokenize(t *testing.T) {
	program := bytes.NewReader([]byte(`
Ook. Ook? Ook? Ook. Ook. Ook.
Ook! Ook! Ook! Ook. Ook. Ook! Ook! Ook?
Ook? Ook! Ook? Ook?
`))

	tokens, err := Tokenize(program)
	assert.NoError(t, err)

	assert.Equal(t, []bf.Token{
		bf.MoveRight, bf.MoveLeft, bf.Increment, bf.Decrement, bf.Output, bf.Input, bf.JumpForward, bf.JumpBackward, bf.Noop,
	}, tokens)

}
