package chip8

// clearDisplay blanks the display and marks it dirty for redraw.
func (c *Chip8) clearDisplay() {
	for i := 0; i < len(c.display); i++ {
		for j := 0; j < len(c.display[i]); j++ {
			c.display[i][j] = false
		}
	}
	c.render = true
}

// draw handles DXYN: it XORs an N-byte sprite from memory (starting at the
// index register) onto the display at (VX, VY), wrapping the start position
// but clipping the sprite at the screen edges, and sets VF if any pixel was
// erased.
func (c *Chip8) draw(X, Y, N uint8) {
	px := c.register[X] % 64
	py := c.register[Y] % 32
	c.register[0xF] = 0
	c.debug.Println("x,y: ", px, py)
	for row := 0; row < int(N); row++ {
		if int(py) >= len(c.display) { // reached bottom edge of screen
			break
		}
		sbyte := c.memory[c.i+uint16(row)]
		c.debug.Printf("%s %x\n", "s: ", sbyte)
		px := c.register[X] % 64
		for bit := 0; bit < 8; bit++ {
			if int(px) >= len(c.display[0]) { // reached right edge of screen
				break
			}
			spritePixel := (sbyte >> (7 - bit)) & 0x01
			displayPixel := c.display[py][px]
			if spritePixel != 0 && displayPixel {
				c.display[py][px] = false
				c.register[0xF] = 1
			} else if spritePixel != 0 {
				c.display[py][px] = true
			}
			px += 1
		}
		py += 1
	}
	c.render = true
}
