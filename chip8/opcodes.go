package chip8

import (
	"fmt"
	"math/rand"
)

func (c *Chip8) Execute(opcode uint16) error {
	// nibbles
	n1 := opcode & 0xF000
	X := uint8((opcode & 0x0F00) >> 8)
	Y := uint8((opcode & 0x00F0) >> 4)
	N := uint8(opcode & 0x000F)
	NN := uint8(opcode & 0x00FF)
	NNN := opcode & 0x0FFF

	c.debug.Printf("%s %x\n", "nibble:", (n1))
	c.debug.Printf("%s %d\n", "X: ", X)
	c.debug.Printf("%s %d\n", "Y: ", Y)
	c.debug.Printf("%s %d\n", "N: ", N)
	c.debug.Printf("%s %x - %d\n", "NN: ", NN, NN)
	c.debug.Printf("%s %x - %d\n", "NNN: ", NNN, NNN)
	switch n1 {

	case 0x0000:
		return c.execute0NNN(opcode, NN)
	case 0x1000: // jump
		c.pc = NNN
	case 0x2000: // call subroutine
		if int(c.sp) >= len(c.stack) {
			return fmt.Errorf("stack overflow")
		}
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = NNN
	case 0x3000: // skip if true
		if c.register[X] == NN {
			c.pc += 2
		}
	case 0x4000: // skip if not true
		if c.register[X] != NN {
			c.pc += 2
		}
	case 0x5000: // skip if true
		if c.register[X] == c.register[Y] {
			c.pc += 2
		}
	case 0x6000:
		// Set
		c.register[X] = NN
	case 0x7000:
		// Add
		c.register[X] += NN
	case 0x8000:
		return c.execute8XY(opcode, X, Y, N)
	case 0x9000: // skip if not true
		if c.register[X] != c.register[Y] {
			c.pc += 2
		}
	case 0xA000:
		c.i = NNN
	case 0xB000:
		// jump to address NNN + V0
		c.pc = NNN + uint16(c.register[0])
	case 0xC000: // random
		c.register[X] = uint8(rand.Uint32()&0xFF) & NN
	case 0xD000:
		c.draw(X, Y, N)
	case 0xE000:
		c.executeE(X, NN)
	case 0xF000:
		c.executeF(X, NN)

	default:
		return fmt.Errorf("unknown opcode: %#04X", opcode)
	}

	return nil
}
