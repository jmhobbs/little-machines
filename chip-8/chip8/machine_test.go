package chip8

import "testing"

/*
	In a two byte op, the values can be
	in the following nibbles:

  +-----+---------------+
  |     | 0 | 1 | 2 | 3 |
	+-----+---------------+
	| nnn |   | n | n | n |
	+-----+---------------+
	|  n  |   |   |   | n |
	+-----+---------------+
	|  x  |   | x |   |   |
	+-----+---------------+
	|  y  |   |   | y |   |
	+-----+---------------+
	| kk  |   |   | k | k |
	+-----+---------------+

*/

func Test_nnn(t *testing.T) {
	var op uint16 = 0b1111_101010101010
	var expected uint16 = 0b0000_101010101010
	actual := nnn(op)

	if expected != actual {
		t.Errorf("mismatch\n  actual: 0b%016b\nexpected: 0b%016b", actual, expected)
	}
}

func Test_n(t *testing.T) {
	var op uint16 = 0b000000000000_1010
	var expected uint8 = 0b0000_1010
	actual := n(op)

	if expected != actual {
		t.Errorf("mismatch\n  actual: 0b%08b\nexpected: 0b%08b", actual, expected)
	}
}

func Test_kk(t *testing.T) {
	var op uint16 = 0b00000000_11101011
	var expected uint8 = 0b11101011
	actual := kk(op)

	if expected != actual {
		t.Errorf("mismatch\n  actual: 0b%08b\nexpected: 0b%08b", actual, expected)
	}
}

func Test_x(t *testing.T) {
	var op uint16 = 0b1111_1010_00000000
	var expected uint8 = 0b0000_1010
	actual := x(op)

	if expected != actual {
		t.Errorf("mismatch\n  actual: 0b%08b\nexpected: 0b%08b", actual, expected)
	}
}

func Test_y(t *testing.T) {
	var op uint16 = 0b11111111_1010_0000
	var expected uint8 = 0b0000_1010
	actual := y(op)

	if expected != actual {
		t.Errorf("mismatch\n  actual: 0b%08b\nexpected: 0b%08b", actual, expected)
	}
}
