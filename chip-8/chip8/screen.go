package chip8

import (
	"fmt"
	"strings"
)

// 1 bit per pixel
// rows are 8 bytes wide
// there are 32 rows
type screen [256]uint8

func (s *screen) Write(sprites []byte, x, y uint8) bool {
	collision := false

	n := x % 8
	x1 := x / 8
	x2 := x1 + 1
	if x2 > 7 { // wrap around
		x2 = 0
	}

	for i := uint8(0); i < uint8(len(sprites)); i++ {
		if !collision && s[(y+i)*8+x]&sprites[i] != 0 {
			collision = true
		}
		
		if n == 0 {
			// byte aligned, we only deal with x1
			s[(y+i)*8+x1] = s[(y+i)*8+x1] ^ sprites[i]
		} else {
			// not byte aligned, we do both
			s[(y+i)*8+x1] = s[(y+i)*8+x1] ^ (sprites[i] >> n)
			s[(y+i)*8+x2] = s[(y+i)*8+x2] ^ (sprites[i] << (8-n))
		}
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
