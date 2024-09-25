package chip8

// CHIP-8 uses a 16 key hex keyboard. 0-F

type Keyboard interface {
	// currently pressed keys
	Pressed() []uint8
	// wait for a key to be pressed
	WaitForPress() uint8
}
