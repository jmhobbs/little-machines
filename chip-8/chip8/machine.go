package chip8

import (
	"context"
	"crypto/rand"
	"fmt"
)

// Based almost entirely on http://devernay.free.fr/hacks/chip8/C8TECH10.HTM

type Machine interface {
	Run(context.Context) error
	Screen() [256]uint8
	Step() error
	State() State
}

func New(pgm []byte, keyboard Keyboard) (Machine, error) {
	if len(pgm) > 3584 { // 4096 - 512
		return nil, fmt.Errorf("program does not fit in memory: %d > 3584", len(pgm))
	}

	m := &machine{
		keyboard: keyboard,
		memory: [4096]byte{
			// font 0-F
			0xF0, 0x90, 0x90, 0x90, 0xF0,
			0x20, 0x60, 0x20, 0x20, 0x70,
			0xF0, 0x10, 0xF0, 0x80, 0xF0,
			0xF0, 0x10, 0xF0, 0x10, 0xF0,
			0x90, 0x90, 0xF0, 0x10, 0x10,
			0xF0, 0x80, 0xF0, 0x10, 0xF0,
			0xF0, 0x80, 0xF0, 0x90, 0xF0,
			0xF0, 0x10, 0x20, 0x40, 0x40,
			0xF0, 0x90, 0xF0, 0x90, 0xF0,
			0xF0, 0x90, 0xF0, 0x10, 0xF0,
			0xF0, 0x90, 0xF0, 0x90, 0x90,
			0xE0, 0x90, 0xE0, 0x90, 0xE0,
			0xF0, 0x80, 0x80, 0x80, 0xF0,
			0xE0, 0x90, 0x90, 0x90, 0xE0,
			0xF0, 0x80, 0xF0, 0x80, 0xF0,
			0xF0, 0x80, 0xF0, 0x80, 0x80,
		},
		PC: 0x0200,
	}

	for i := 0; i < len(pgm); i++ {
		m.memory[i+0x200] = pgm[i]
	}

	return m, nil
}

type machine struct {
	keyboard Keyboard
	memory   [4096]uint8
	screen   screen     // Screen buffer 64x32 monopixel
	stack    [16]uint16 // Stack
	PC       uint16     // Program counter
	SP       uint8      // Stack pointer
	V        [16]uint8  // General purpose registers 0-F
	I        uint16     // Memory address register (12 bit addresses)
	DT       uint8      // Delay Timer, decremented at 60hz until 0
	ST       uint8      // Sound Timer, decremented at 60hz until 0
}

type State struct {
	stack [16]uint16 // Stack
	PC    uint16     // Program counter
	SP    uint8      // Stack pointer
	V     [16]uint8  // General purpose registers 0-F
	I     uint16     // Memory address register (12 bit addresses)
	DT    uint8      // Delay Timer, decremented at 60hz until 0
	ST    uint8      // Sound Timer, decremented at 60hz until 0
}

func (m *machine) Screen() [256]uint8 {
	return m.screen.Bytes()
}

func (m *machine) State() State {
	return State{
		stack: m.stack,
		PC:    m.PC,
		SP:    m.SP,
		V:     m.V,
		I:     m.I,
		DT:    m.DT,
		ST:    m.ST,
	}
}

func (m *machine) Run(ctx context.Context) error {
	// todo: DT & ST decay at 60hz
	// todo: run at 500hz
	var err error
	for {
		select {
		case <-ctx.Done():
			return nil // returning not to leak the goroutine
		default:
			err = m.Step()
			if err != nil {
				return err
			}
		}
	}
}

const (
	// fully specified
	CLS uint16 = 0x00E0
	RET uint16 = 0x00EE
	// masks
	JP    uint16 = 0x1
	CALL  uint16 = 0x2
	SE    uint16 = 0x3
	SNE   uint16 = 0x4
	SE_V  uint16 = 0x5
	LD    uint16 = 0x6
	ADD   uint16 = 0x7
	SNE_V uint16 = 0x9
	LD_I  uint16 = 0xA
	JP_V0 uint16 = 0xB
	RND   uint16 = 0xC
	DRW   uint16 = 0xD
	// Lower nibble, 8__N
	LD_V   uint8 = 0x00
	OR_V   uint8 = 0x01
	AND_V  uint8 = 0x02
	XOR_V  uint8 = 0x03
	ADD_V  uint8 = 0x04
	SUB_V  uint8 = 0x05
	SHR_V  uint8 = 0x06
	SUBN_V uint8 = 0x07
	SHL_V  uint8 = 0x08
	// Lower byte, E_NN
	SKP  uint8 = 0x9E
	SKNP uint8 = 0xA1
	// Lower byte, F_NN
	LD_V_DT  uint8 = 0x07
	LD_K     uint8 = 0x0A
	LD_DT_V  uint8 = 0x15
	LD_ST_V  uint8 = 0x18
	ADD_I_V  uint8 = 0x1E
	LD_FONT  uint8 = 0x29
	LD_BCD   uint8 = 0x33
	LD_V_MEM uint8 = 0x55
	LD_MEM_V uint8 = 0x65
)

func (m *machine) Step() error {
	advPC := true
	op := (uint16(m.memory[m.PC]) << 8) | uint16(m.memory[m.PC+1])

	// todo: there must be some bit matching that is better than this nested mess
	upperNibble := op >> 12
	lowerByte := kk(op)

	if op == CLS { // 00E0
		m.screen.Clear()
	} else if op == RET { // 00EE
		// todo: stack underflow
		m.SP--
		m.PC = m.stack[m.SP]
		advPC = false
	} else if upperNibble == JP { // 1nnn
		m.PC = nnn(op)
		advPC = false
	} else if upperNibble == CALL { // 2nnn
		// todo: stack overflow
		m.stack[m.SP] = m.PC
		m.SP++
		m.PC = nnn(op)
		advPC = false
	} else if upperNibble == SE { // 3xkk
		if m.V[x(op)] == lowerByte {
			m.PC += 2
		}
	} else if upperNibble == SNE { // 4xkk
		if m.V[x(op)] != lowerByte {
			m.PC += 2
		}
	} else if upperNibble == SE_V && lowerByte&0xF == 0 { // 5xy0
		if m.V[x(op)] == m.V[y(op)] {
			m.PC += 2
		}
	} else if upperNibble == LD { // 6xkk
		m.V[x(op)] = kk(op)
	} else if upperNibble == ADD { // 7xkk
		m.V[x(op)] += kk(op)
	} else if upperNibble == 0x8 && lowerByte == LD_V { // 8xy0
		m.V[x(op)] = m.V[y(op)]
	} else if upperNibble == 0x8 && lowerByte == OR_V { // 8xy1
		m.V[x(op)] = m.V[x(op)] | m.V[y(op)]
	} else if upperNibble == 0x8 && lowerByte == AND_V { // 8xy2
		m.V[x(op)] = m.V[x(op)] & m.V[y(op)]
	} else if upperNibble == 0x8 && lowerByte == XOR_V { // 8xy3
		m.V[x(op)] = m.V[x(op)] ^ m.V[y(op)]
	} else if upperNibble == 0x8 && lowerByte == ADD_V { // 8xy4
		v := uint16(m.V[x(op)]) + uint16(m.V[y(op)])
		if v > 255 {
			m.V[0xF] = 1
		} else {
			m.V[0xF] = 0
		}
		m.V[x(op)] = uint8(v)
	} else if upperNibble == 0x8 && lowerByte == SUB_V { // 8xy5
		X := m.V[x(op)]
		Y := m.V[y(op)]
		if X > Y {
			m.V[0xF] = 1
		} else {
			m.V[0xF] = 0
		}
		m.V[x(op)] = X - Y
	} else if upperNibble == 0x8 && lowerByte == SHR_V { // 8xy6
		X := m.V[x(op)]
		if X&0x01 == 0x01 {
			m.V[0xF] = 1
		} else {
			m.V[0xF] = 0
		}
		// todo: rounding semantics?
		m.V[x(op)] = X / 2
	} else if upperNibble == 0x8 && lowerByte == SUBN_V { // 8xy7
		X := m.V[x(op)]
		Y := m.V[y(op)]
		if Y > X {
			m.V[0xF] = 1
		} else {
			m.V[0xF] = 0
		}
		m.V[x(op)] = Y - X
	} else if upperNibble == 0x8 && lowerByte == SHL_V { // 8xyE
		X := m.V[x(op)]
		if X&0b10000000 == 0b10000000 {
			m.V[0xF] = 1
		} else {
			m.V[0xF] = 0
		}
		m.V[x(op)] = X * 2
	} else if upperNibble == SNE_V && lowerByte&0x0F == 0 { // 9xy0
		if m.V[x(op)] != m.V[y(op)] {
			m.PC += 2
		}
	} else if upperNibble == LD_I { // Annn
		m.I = nnn(op)
	} else if upperNibble == JP_V0 { // Bnnn
		m.PC = uint16(m.V[0]) + nnn(op)
		advPC = false
	} else if upperNibble == RND { // Cnnn
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			return err
		}
		m.V[x(op)] = b[0] & kk(op)
	} else if upperNibble == DRW { // Dxyn
		m.screen.Write(m.memory[m.I:m.I+uint16(n(op))], m.V[x(op)], m.V[y(op)])
	} else if upperNibble == 0xE && lowerByte == SKP { // Ex9E
		keys := m.keyboard.Pressed()
		for _, k := range keys {
			if k == m.V[x(op)] {
				m.PC += 2
				break
			}
		}
	} else if upperNibble == 0xE && lowerByte == SKNP { // ExA1
		keys := m.keyboard.Pressed()
		isPressed := false
		for _, k := range keys {
			if k == m.V[x(op)] {
				isPressed = true
				break
			}
		}
		if !isPressed {
			m.PC += 2
		}
	} else if upperNibble == 0xF && lowerByte == LD_V_DT { // Fx07
		m.V[x(op)] = m.DT
	} else if upperNibble == 0xF && lowerByte == LD_K { // Fx0A
		m.V[x(op)] = m.keyboard.WaitForPress()
	} else if upperNibble == 0xF && lowerByte == LD_DT_V { // Fx15
		m.DT = m.V[x(0)]
	} else if upperNibble == 0xF && lowerByte == LD_ST_V { // Fx18
		m.ST = m.V[x(0)]
	} else if upperNibble == 0xF && lowerByte == ADD_I_V { // Fx1E
		m.I = m.I + uint16(m.V[x(op)])
	} else if upperNibble == 0xF && lowerByte == LD_FONT { // Fx29
		m.I = uint16(m.V[x(op)] & 0x0F * 8)
	} else if upperNibble == 0xF && lowerByte == LD_BCD { // Fx33
		m.memory[m.I] = m.V[x(op)] / 100
		m.memory[m.I+1] = (m.V[x(op)] % 100) / 10
		m.memory[m.I+2] = ((m.V[x(op)] % 100) % 10) / 1
	} else if upperNibble == 0xF && lowerByte == LD_V_MEM { // Fx55
		for x := uint16(0); x < 16; x++ {
			m.memory[m.I+x] = m.V[x]
		}
	} else if upperNibble == 0xF && lowerByte == LD_MEM_V { // Fx65
		for x := uint16(0); x < 16; x++ {
			m.V[x] = m.memory[m.I+x]
		}
	} else {
		return fmt.Errorf("unknown opcode: 0x%04x", op)
	}

	if advPC {
		m.PC += 2
	}

	return nil
}

// get the n byte (value)
func n(op uint16) uint8 {
	return uint8(op) & 0x0F
}

// get the kk byte (value)
func kk(op uint16) uint8 {
	return uint8(op)
}

// get the 0nnn byte (address)
func nnn(op uint16) uint16 {
	return op & 0x0FFF
}

// get the x nibble (register number)
func x(op uint16) uint8 {
	return uint8(op>>8) & 0x0F
}

// get the y nibble from (secondary register number)
func y(op uint16) uint8 {
	return uint8(op>>4) & 0x0F
}
