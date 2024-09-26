package asm

import "fmt"

/**
 * Read a register value from a string.
 * Returns an error if out of bounds or improperly formed.
 * Expects uppercase.
 */
func register(v string) (uint8, error) {
	var n uint8
	_, err := fmt.Sscanf(v, "V%x", &n)
	if err != nil || n > 15 {
		return 0, fmt.Errorf("invalid register %q", v)
	}
	return n, err
}

/**
 * Read two register values from strings.
 * Returns an error if out of bounds or improperly formed.
 * Expects uppercase.
 */
func registers(a, b string) (uint8, uint8, error) {
	ar, err := register(a)
	if err != nil {
		return 0, 0, err
	}
	br, err := register(b)
	return ar, br, err
}

/**
 * Read a register value and a byte value from strings.
 * Returns an error if out of bounds or improperly formed.
 * Expects uppercase.
 */
func registerValue(a, b string) (uint8, uint8, error) {
	ar, err := register(a)
	if err != nil {
		return 0, 0, err
	}
	var value uint8
	_, err = fmt.Sscanf(b, "%x", &value)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid value %q", b)
	}
	return ar, value, nil
}

type Operand uint8

var (
	REGISTER Operand = 0
	BYTE     Operand = 1
)

func operandType(s string) Operand {
	if s[0] == 'V' {
		return REGISTER
	}
	return BYTE
}
