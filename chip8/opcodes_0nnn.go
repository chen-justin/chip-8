package chip8

import "fmt"

// execute0NNN handles the 0x0NNN opcode family: clearing the display (00E0)
// and returning from a subroutine (00EE).
func (c *Chip8) execute0NNN(opcode uint16, NN uint8) error {
	switch NN {
	case 0xE0: // clear Screen
		c.clearDisplay()
	case 0xEE: // return subroutine
		if c.sp == 0 {
			return fmt.Errorf("returned with empty stack")
		}
		c.sp -= 1
		c.pc = c.stack[c.sp]
	default:
		return fmt.Errorf("unknown 0x0NNN opcode: %#04X", opcode)
	}
	return nil
}
