package chip8

import "fmt"

// execute8XY handles the 0x8XY? opcode family: register-to-register ALU
// operations (set, bitwise ops, add/subtract with carry/borrow, and shift).
func (c *Chip8) execute8XY(opcode uint16, X, Y, N uint8) error {
	switch N {
	case 0x0000:
		//set
		c.register[X] = c.register[Y]
	case 0x0001:
		// binary or
		c.register[X] = c.register[X] | c.register[Y]
	case 0x0002:
		// binary and
		c.register[X] = c.register[X] & c.register[Y]
	case 0x0003:
		// logic xor
		c.register[X] = c.register[X] ^ c.register[Y]
	case 0x0004:
		// add, VF = carry
		sum := uint16(c.register[X]) + uint16(c.register[Y])
		c.register[X] = uint8(sum)
		if sum > 0xFF {
			c.register[0xF] = 1
		} else {
			c.register[0xF] = 0
		}
	case 0x0005:
		// subtract VX - VY, VF = NOT borrow
		borrow := c.register[X] < c.register[Y]
		c.register[X] = c.register[X] - c.register[Y]
		if borrow {
			c.register[0xF] = 0
		} else {
			c.register[0xF] = 1
		}
	case 0x0007:
		// subtract VY - VX, VF = NOT borrow
		borrow := c.register[Y] < c.register[X]
		c.register[X] = c.register[Y] - c.register[X]
		if borrow {
			c.register[0xF] = 0
		} else {
			c.register[0xF] = 1
		}
	case 0x0006:
		// shift
		c.register[X] = c.register[Y]
		carry := c.register[X] & 0x01
		c.register[X] = c.register[X] >> 1
		c.register[15] = carry
	case 0x000E:
		c.register[X] = c.register[Y]
		carry := (c.register[X] & 0x80) >> 7
		c.register[X] = c.register[X] << 1
		c.register[0xF] = carry
	default:
		return fmt.Errorf("unknown 0x8XY%X opcode: %#04X", N, opcode)
	}
	return nil
}
