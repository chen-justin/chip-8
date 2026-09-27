package chip8

// executeE handles the 0xE... opcode family: skip-if-key-pressed checks.
func (c *Chip8) executeE(X uint8, NN uint8) {
	switch NN {
	case 0x9E:
		if c.key[X] {
			c.pc += 2
		}
	case 0xA1:
		if !c.key[X] {
			c.pc += 2
		}
	}
}
