package bf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PointerAdvance(t *testing.T) {
	i := &interpreter{
		m:       NewMachine(300),
		program: []Token("++>[++[>+<]--<]"),
		ptr:     0,
	}

	i.Step()
	assert.Equal(t, uint(1), i.ptr)
	i.Step()
	assert.Equal(t, uint(2), i.ptr)

}

func Test_IncrementAndDecrement(t *testing.T) {
	i := &interpreter{
		m:       NewMachine(300),
		program: []Token("++-"),
		ptr:     0,
	}

	assert.Equal(t, byte(0), i.m.Peek())
	i.Step()
	i.Step()
	assert.Equal(t, byte(2), i.m.Peek())
	i.Step()
	assert.Equal(t, byte(1), i.m.Peek())
}

func Test_MemoryPointerIncrementAndDecrement(t *testing.T) {
	i := &interpreter{
		m:       NewMachine(300),
		program: []Token("+>++<"),
		ptr:     0,
	}

	i.Run()
	assert.Equal(t, []byte{1, 2, 0}, i.m.Dump(0, 3))
}

func Test_FindLoopEnd(t *testing.T) {
	i := &interpreter{
		m:       NewMachine(300),
		program: []Token("++>[++[>+<]--<]"),
		ptr:     3,
	}

	actual, err := i.findLoopEnd()
	if assert.NoError(t, err) {
		assert.Equal(t, uint(14), actual)
		assert.Equal(t, JumpBackward, i.program[actual])
	}

	i.ptr = 6
	actual, err = i.findLoopEnd()
	if assert.NoError(t, err) {
		assert.Equal(t, uint(10), actual)
		assert.Equal(t, JumpBackward, i.program[actual])
	}
}

func Test_FindLoopStart(t *testing.T) {
	i := &interpreter{
		m:       NewMachine(300),
		program: []Token("++>[++[>+<]--<]"),
		ptr:     14,
	}

	actual, err := i.findLoopStart()
	if assert.NoError(t, err) {
		assert.Equal(t, uint(3), actual)
		assert.Equal(t, JumpForward, i.program[actual])
	}

	i.ptr = 10
	actual, err = i.findLoopStart()
	if assert.NoError(t, err) {
		assert.Equal(t, uint(6), actual)
		assert.Equal(t, JumpForward, i.program[actual])
	}
}
