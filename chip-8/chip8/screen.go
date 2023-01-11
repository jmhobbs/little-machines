package chip8

import (
	"fmt"
	"strings"
)

// 1 bit per pixel
// rows are 8 bytes wide
// there are 32 rows
type screen [256]uint8

// todo: wrapping
func (s *screen) Write(sprites []byte, x, y uint8) bool {
	collision := false

	for i := uint8(0); i < uint8(len(sprites)); i++ {
		if !collision && s[(y+i)*8+x]&sprites[i] != 0 {
			collision = true
		}
		s[(y+i)*8+x] = s[(y+i)*8+x] ^ sprites[i]
	}

	return collision
}

func (s *screen) Clear() {
	s = &screen{}
}

func (s *screen) String() string {
	var b strings.Builder

	for y := 0; y < 32; y++ {
		for x := 0; x < 8; x++ {
			b.WriteString(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%08b", s[y*8+x]), "1", "▣"), "0", "▢"))
		}
		b.WriteRune('\n')
	}

	return b.String()
}
