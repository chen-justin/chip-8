package chip8

// executeF handles the 0xF... opcode family: timers, index register
// arithmetic, the font-character lookup, BCD conversion, blocking key
// input, and bulk register/memory transfers.
func (c *Chip8) executeF(X uint8, NN uint8) {
	switch NN {
	// timers
	case 0x07:
		c.register[X] = c.dt
	case 0x15:
		c.dt = c.register[X]
	case 0x18:
		c.st = c.register[X]
	case 0x1E: // add to index
		c.i += uint16(c.register[X])
	case 0x0A: // get key: block until some key is pressed
		pressed := false
		for k := uint8(0); k < uint8(len(c.key)); k++ {
			if c.key[k] {
				c.register[X] = k
				pressed = true
				break
			}
		}
		if !pressed {
			c.pc -= 2 // re-fetch this same instruction next cycle
		}
	case 0x29: //font character
		c.i = 0x50 + uint16(c.register[X])*5 // every font character is 5 bytes
	case 0x33: //binary-coded decimal conversion
		temp := c.register[X]
		i := 2
		for i >= 0 {
			digit := temp % 10
			c.memory[c.i+uint16(i)] = digit
			temp /= 10
			i--
		}
	case 0x55:
		for i := 0; i <= int(X); i += 1 {
			c.memory[c.i+uint16(i)] = c.register[i]
		}
	case 0x65:
		for i := 0; i <= int(X); i += 1 {
			c.register[i] = c.memory[c.i+uint16(i)]
		}
	}
}
