package bf

import (
	"fmt"
	"io"
)

func read() byte {
	fmt.Print("input: ")
	var r rune
	_, err := fmt.Scanf("%%c", &r)
	if err != nil {
		panic(err)
	}
	return byte(r)
}

func NewInterpreter(size uint, program []Token) Interpreter {
	return &interpreter{
		m:       NewMachine(size),
		program: program,
		ptr:     0,
	}
}

type Interpreter interface {
	Run()
	Step()
	Dump(uint, uint) []byte
}

type interpreter struct {
	m       Machine
	program []Token
	ptr     uint
}

func (i *interpreter) Run() {
	length := uint(len(i.program))
	for i.ptr < length {
		i.Step()
	}
}

func (i *interpreter) Dump(start, end uint) []byte {
	return i.m.Dump(start, end)
}

func (i *interpreter) Step() {
	switch i.program[i.ptr] {
	case MoveRight:
		i.m.PtrIncr()
	case MoveLeft:
		i.m.PtrDecr()
	case Increment:
		i.m.Incr()
	case Decrement:
		i.m.Decr()
	case Output:
		fmt.Print(string(i.m.Peek()))
	case Input:
		i.m.Put(read())
	case JumpForward:
		if i.m.Peek() == 0 {
			n, err := i.findLoopEnd()
			if err != nil {
				panic(err)
			}
			i.ptr = n
		}
	case JumpBackward:
		// todo: optimize the check here instead
		n, err := i.findLoopStart()
		if err != nil {
			panic(err)
		}
		i.ptr = n - 1
	default:
		// noop
	}
	i.ptr++
}

func (i *interpreter) findLoopEnd() (uint, error) {
	length := uint(len(i.program))
	openLoops := 0
	for j := i.ptr + 1; j < length; j++ {
		if i.program[j] == JumpForward {
			openLoops++
		} else if i.program[j] == JumpBackward {
			if openLoops == 0 {
				return j, nil
			}
			openLoops--
		}
	}
	return length - 1, io.EOF
}

func (i *interpreter) findLoopStart() (uint, error) {
	length := uint(len(i.program))
	openLoops := 0
	for j := i.ptr - 1; j < length; j-- {
		if i.program[j] == JumpBackward {
			openLoops++
		} else if i.program[j] == JumpForward {
			if openLoops == 0 {
				return j, nil
			}
			openLoops--
		}
	}
	return length - 1, io.EOF
}
