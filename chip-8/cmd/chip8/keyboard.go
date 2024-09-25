package main

import "github.com/jmhobbs/little-machines/chip-8/chip8"

var _ chip8.Keyboard = (*Keyboard)(nil)

type Keyboard struct{}

func (k *Keyboard) Pressed() []uint8 {
	return []uint8{}
}

func (k *Keyboard) WaitForPress() uint8 {
	return 0xF
}
