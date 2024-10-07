package asm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encoder(t *testing.T) {

	tests := []struct {
		instruction string
		operands    []string
		expected    []byte
	}{
		{
			"SYS", []string{"0X234"}, []byte{0x02, 0x34},
		},
		{
			"CLS", []string{}, []byte{0x00, 0xE0},
		},
		{
			"RET", []string{}, []byte{0x00, 0xEE},
		},
		{
			"JUMP", []string{"0X234"}, []byte{0x12, 0x34},
		},
		{
			"CALL", []string{"0X234"}, []byte{0x22, 0x34},
		},
		{
			"SKIP", []string{"V1", "0X34"}, []byte{0x31, 0x34},
		},
		{
			"SKIPN", []string{"V2", "0X43"}, []byte{0x42, 0x43},
		},
		{
			"LOAD", []string{"VA", "0X11"}, []byte{0x6A, 0x11},
		},
		{
			"ADD", []string{"V7", "0X51"}, []byte{0x77, 0x51},
		},
		{
			"LOAD", []string{"V1", "V2"}, []byte{0x81, 0x20},
		},
		{
			"OR", []string{"V1", "V2"}, []byte{0x81, 0x21},
		},
		{
			"AND", []string{"V5", "V4"}, []byte{0x85, 0x42},
		},
		{
			"XOR", []string{"V5", "V2"}, []byte{0x85, 0x23},
		},
		{
			"ADD", []string{"V8", "V0"}, []byte{0x88, 0x04},
		},
		{
			"SUB", []string{"V9", "V1"}, []byte{0x89, 0x15},
		},
		{
			"SHIFTR", []string{"V5"}, []byte{0x85, 0x06},
		},
		{
			"SUBN", []string{"V5", "V4"}, []byte{0x85, 0x47},
		},
		{
			"SHIFTL", []string{"V5"}, []byte{0x85, 0x0E},
		},
		{
			"SKIPN", []string{"V6", "VC"}, []byte{0x96, 0xC0},
		},
		{
			"LOAD", []string{"I", "0X123"}, []byte{0xA1, 0x23},
		},
		{
			"JUMP", []string{"V0", "0X333"}, []byte{0xB3, 0x33},
		},
		{
			"RAND", []string{"V1", "0X33"}, []byte{0xC1, 0x33},
		},
		{
			"DRAW", []string{"V1", "V2", "0X03"}, []byte{0xD1, 0x23},
		},
		{
			"SKIPP", []string{"V6"}, []byte{0xE6, 0x9E},
		},
		{
			"SKIPNP", []string{"V6"}, []byte{0xE6, 0xA1},
		},
		{
			"LOAD", []string{"V9", "DT"}, []byte{0xF9, 0x07},
		},
		{
			"LOAD", []string{"VD", "KEY"}, []byte{0xFD, 0x0A},
		},
		{
			"LOAD", []string{"DT", "V9"}, []byte{0xF9, 0x15},
		},
		{
			"LOAD", []string{"ST", "VE"}, []byte{0xFE, 0x18},
		},
		{
			"ADD", []string{"I", "V3"}, []byte{0xF3, 0x1E},
		},
		{
			"LOAD", []string{"FONT", "V3"}, []byte{0xF3, 0x29},
		},
		{
			"LOAD", []string{"BCD", "V7"}, []byte{0xF7, 0x33},
		},
		{
			"LOAD", []string{"I", "V8"}, []byte{0xF8, 0x55},
		},
		{
			"LOAD", []string{"V9", "I"}, []byte{0xF9, 0x65},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %v", test.instruction, test.operands), func(t *testing.T) {
			actual, err := Encode(test.instruction, test.operands)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}

}
