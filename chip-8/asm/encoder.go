package asm

import (
	"fmt"
	"slices"
)

type encoder func([]string) ([]byte, error)

type encoding struct {
	operands []Operand
	encoder  encoder
}

var instructionAliases = map[string]string{
	"CLS": "CLEAR",
	"RET": "RETURN",
	"JMP": "JUMP",
	"SE":  "SKIP",
	"SNE": "SKIPN",
	"LD":  "LOAD",
	"SHR": "SHIFTR",
	"SHL": "SHIFTL",
}

var instructions = map[string][]encoding{
	"CLEAR": {
		{
			encoder: func([]string) ([]byte, error) {
				return []byte{0x00, 0xE0}, nil
			},
		},
	},
	"RETURN": {
		{
			encoder: func([]string) ([]byte, error) {
				return []byte{0x00, 0xEE}, nil
			},
		},
	},
	"JUMP": {
		addressEncoding(0x10),
		// V0, addr - Bnnn
	},
	"CALL": {
		addressEncoding(0x20),
	},
	"SKIP": {
		registerByteEncoding(0x30),
		twoRegisterEncoding(0x50, 0x00),
	},
	"SKIPN": {
		registerByteEncoding(0x40),
		twoRegisterEncoding(0x90, 0x00),
	},
	"LOAD": {
		registerByteEncoding(0x60),
		twoRegisterEncoding(0x80, 0x00),
		{
			operands: []Operand{I_REGISTER, ADDRESS},
			encoder:  addressEncoder(0xA0),
		},
		{
			operands: []Operand{REGISTER, DELAY_TIMER},
			encoder:  oneRegisterEncoder(0xF0, 0x07),
		},
		{
			operands: []Operand{REGISTER, KEYPRESS},
			encoder:  oneRegisterEncoder(0xF0, 0x0A),
		},
		{
			operands: []Operand{DELAY_TIMER, REGISTER},
			encoder:  oneRegisterEncoderSkip(0xF0, 0x15),
		},
		{
			operands: []Operand{SOUND_TIMER, REGISTER},
			encoder:  oneRegisterEncoderSkip(0xF0, 0x18),
		},
		{
			operands: []Operand{SPRITE, REGISTER},
			encoder:  oneRegisterEncoderSkip(0xF0, 0x29),
		},
		{
			operands: []Operand{BCD, REGISTER},
			encoder:  oneRegisterEncoderSkip(0xF0, 0x33),
		},
		{
			operands: []Operand{I_REGISTER, REGISTER},
			encoder:  oneRegisterEncoderSkip(0xF0, 0x55),
		},
		{
			operands: []Operand{REGISTER, I_REGISTER},
			encoder:  oneRegisterEncoder(0xF0, 0x65),
		},
	},
	"ADD": {
		registerByteEncoding(0x70),
		twoRegisterEncoding(0x80, 0x04),
		{
			operands: []Operand{I_REGISTER, REGISTER},
			encoder:  oneRegisterEncoderSkip(0xF0, 0x1E),
		},
	},
	"OR": {
		twoRegisterEncoding(0x80, 0x01),
	},
	"AND": {
		twoRegisterEncoding(0x80, 0x02),
	},
	"XOR": {
		twoRegisterEncoding(0x80, 0x03),
	},
	"SUB": {
		twoRegisterEncoding(0x80, 0x05),
	},
	"SHIFTR": {
		oneRegisterEncoding(0x80, 0x06),
	},
	"SUBN": {
		twoRegisterEncoding(0x80, 0x07),
	},
	"SHIFTL": {
		oneRegisterEncoding(0x80, 0x0E),
	},
	"RAND": {
		registerByteEncoding(0xC0),
	},
	"DRAW": {
		twoRegisterNibbleEncoding(0xD0),
	},
	"SKIPP": {
		oneRegisterEncoding(0xE0, 0x9E),
	},
	"SKIPNP": {
		oneRegisterEncoding(0xE0, 0xA1),
	},
}

/**
 * Validate and encode an instruction and operands into opcodes.
 */
func Encode(instruction string, operands []string) ([]byte, error) {
	if alias, ok := instructionAliases[instruction]; ok {
		instruction = alias
	}

	op, ok := instructions[instruction]
	if !ok {
		return nil, fmt.Errorf("unknown instruction %q", instruction)
	}

	operandTypes := []Operand{}
	for _, operand := range operands {
		operandTypes = append(operandTypes, operandType(operand))
	}

	for _, encoding := range op {
		if slices.Equal(operandTypes, encoding.operands) {
			return encoding.encoder(operands)
		}
	}

	return nil, fmt.Errorf("invalid operands for %q", instruction)
}

/*
Create an encoder for the one register pattern, with a prefix and a suffix.
The upper bits of the prefix will be used, and the lower bits of the suffix
The encoder will parse registers and return 0xPXYS
*/
func oneRegisterEncoder(prefix, suffix uint8) encoder {
	return func(operands []string) ([]byte, error) {
		target, err := register(operands[0])
		if err != nil {
			return nil, err
		}
		return []byte{(prefix & 0xF0) | target, suffix & 0x0F}, nil
	}
}

func oneRegisterEncoderSkip(prefix, suffix uint8) encoder {
	return func(operands []string) ([]byte, error) {
		target, err := register(operands[1])
		if err != nil {
			return nil, err
		}
		return []byte{(prefix & 0xF0) | target, suffix & 0x0F}, nil
	}
}

func oneRegisterEncoding(prefix, suffix uint8) encoding {
	return encoding{
		operands: []Operand{REGISTER},
		encoder:  oneRegisterEncoder(prefix, suffix),
	}
}

/*
Create an encoder for the two register pattern, with a prefix and a suffix.
The upper bits of the prefix will be used, and the lower bits of the suffix
The encoder will parse registers and return 0xPXYS
*/
func twoRegisterEncoder(prefix, suffix uint8) encoder {
	return func(operands []string) ([]byte, error) {
		dest, src, err := registers(operands[0], operands[1])
		if err != nil {
			return nil, err
		}
		return []byte{(prefix & 0xF0) | dest, src<<4 | (suffix & 0x0F)}, nil
	}
}

func twoRegisterEncoding(prefix, suffix uint8) encoding {
	return encoding{
		operands: []Operand{REGISTER, REGISTER},
		encoder:  twoRegisterEncoder(prefix, suffix),
	}
}

/*
Create an encoder for the register and byte pattern, with a prefix.
The upper bits of the prefix will be used.
The encoder will parse a register and a byte, and return 0xPXNN
*/
func registerByteEncoder(prefix uint8) encoder {
	return func(operands []string) ([]byte, error) {
		dest, value, err := registerValue(operands[0], operands[1])
		if err != nil {
			return nil, err
		}
		return []byte{(prefix & 0xF0) | dest, value}, nil
	}
}

func registerByteEncoding(prefix uint8) encoding {
	return encoding{
		operands: []Operand{REGISTER, BYTE},
		encoder:  registerByteEncoder(prefix),
	}
}

/*
Create an encoder for the address pattern, with a prefix.
The upper bits of the prefix will be used.
The encoder will parse an address and return 0xPNNN
*/
func addressEncoder(prefix uint8) encoder {
	return func(operands []string) ([]byte, error) {
		target, err := address(operands[0])
		if err != nil {
			return nil, err
		}
		return []byte{prefix | uint8(target>>8), uint8(target & 0xFF)}, nil
	}
}

func addressEncoding(prefix uint8) encoding {
	return encoding{
		operands: []Operand{ADDRESS},
		encoder:  addressEncoder(prefix),
	}
}

func twoRegisterNibbleEncoder(prefix uint8) encoder {
	return func(operands []string) ([]byte, error) {
		dest, src, err := registers(operands[0], operands[1])
		if err != nil {
			return nil, err
		}
		nibble, err := value(operands[2])
		if err != nil {
			return nil, err
		}
		return []byte{(prefix & 0xF0) | dest, src<<4 | (nibble & 0x0F)}, nil
	}
}

func twoRegisterNibbleEncoding(prefix uint8) encoding {
	return encoding{
		operands: []Operand{REGISTER, REGISTER, BYTE},
		encoder:  twoRegisterNibbleEncoder(prefix),
	}
}
