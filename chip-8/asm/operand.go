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
	v, err := value(b)
	return ar, v, err
}

/**
 * Read a byte value from a hex encoded string.
 * Returns an error if improperly formed.
 */
func value(a string) (uint8, error) {
	var v uint8
	_, err := fmt.Sscanf(a, "%x", &v)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q", a)
	}
	return v, nil
}

/**
 * Read an address value from a string.
 * Returns an error if out of bounds or improperly formed.
 */
func address(v string) (uint16, error) {
	var n uint16
	_, err := fmt.Sscanf(v, "%x", &n)
	if err != nil || n > 4095 {
		return 0, fmt.Errorf("invalid address %q", v)
	}
	return n, err
}

type Operand uint8

var (
	REGISTER    Operand = 0
	BYTE        Operand = 1
	ADDRESS     Operand = 2
	DELAY_TIMER Operand = 3
	SOUND_TIMER Operand = 4
	I_REGISTER  Operand = 5
	BCD         Operand = 6
	KEYPRESS    Operand = 7
	SPRITE      Operand = 8
)

func operandType(s string) Operand {
	if s[0] == 'V' {
		return REGISTER
	}

	switch s {
	case "DT":
		return DELAY_TIMER
	case "ST":
		return SOUND_TIMER
	case "I":
		return I_REGISTER
	case "B":
		return BCD
	case "K":
		return KEYPRESS
	case "F": // TODO: this will collide with 16 -> F
		return SPRITE
	}

	if len(s) == 3 {
		return ADDRESS
	}
	return BYTE
}
