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

func NewInterpreter(size uint, program []byte) Interpreter {
	return &interpreter{
		m:       NewMachine(300),
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
	program []byte
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
	case '>':
		i.m.PtrIncr()
	case '<':
		i.m.PtrDecr()
	case '+':
		i.m.Incr()
	case '-':
		i.m.Decr()
	case '.':
		fmt.Print(string(i.m.Peek()))
	case ',':
		i.m.Put(read())
	case '[':
		if i.m.Peek() == 0 {
			n, err := i.findLoopEnd()
			if err != nil {
				panic(err)
			}
			i.ptr = n
		}
	case ']':
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
		if i.program[j] == '[' {
			openLoops++
		} else if i.program[j] == ']' {
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
		if i.program[j] == ']' {
			openLoops++
		} else if i.program[j] == '[' {
			if openLoops == 0 {
				return j, nil
			}
			openLoops--
		}
	}
	return length - 1, io.EOF
}
