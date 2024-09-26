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

var instructions = map[string][]encoding{
	"CLEAR": {
		{
			operands: []Operand{},
			encoder: func([]string) ([]byte, error) {
				return []byte{0x00, 0xE0}, nil
			},
		},
	},
	"LOAD": {
		registerByteEncoding(0x60),
		twoRegisterEncoding(0x80, 0x00),
	},
	"ADD": {
		registerByteEncoding(0x70),
		twoRegisterEncoding(0x80, 0x04),
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
	"SHR": {
		twoRegisterEncoding(0x80, 0x06),
	},
	"SHL": {
		twoRegisterEncoding(0x80, 0x0E),
	},
}

/**
 * Validate and encode an instruction and operands into opcodes.
 */
func Encode(instruction string, operands []string) ([]byte, error) {
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
Create an encoder for the two register pattern, with a prefix and a suffix.
The upper bits of the prefix will be used, and the lower bits of the suffix
The encoder will parse registers and return 0xPX 0xYS
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

/*
Create an encoder for the register and byte pattern, with a prefix.
The upper bits of the prefix will be used.
The encoder will parse a register and a byte, and return 0xPX 0xNN
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

func twoRegisterEncoding(prefix, suffix uint8) encoding {
	return encoding{
		operands: []Operand{REGISTER, REGISTER},
		encoder:  twoRegisterEncoder(prefix, suffix),
	}
}

func registerByteEncoding(prefix uint8) encoding {
	return encoding{
		operands: []Operand{REGISTER, BYTE},
		encoder:  registerByteEncoder(prefix),
	}
}
