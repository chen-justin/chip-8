package chip8

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 struct {
	memory   [4096]byte
	display  [32][64]bool
	pc       uint16 //program counter
	i        uint16 //index register
	stack    [16]uint16
	sp       uint16 //stack pointer
	dt       uint8  //delay timer
	st       uint8  //sound timer
	register [16]uint8
	key      [16]bool //keydown
	render   bool
	ips      int //instructions per second
	debug    *log.Logger
}

func Init() Chip8 {
	instance := Chip8{
		pc:     0x200,
		render: true,
		ips:    700,
		debug:  log.New(io.Discard, "", 0), // silent until SetDebug(true)
	}
	buffer := 0x50
	for i := buffer; i < len(fontSet)+buffer; i++ {
		instance.memory[i] = fontSet[i-buffer]
	}
	return instance
}

// SetDebug turns verbose per-instruction tracing to stdout on or off.
func (c *Chip8) SetDebug(on bool) {
	if on {
		c.debug = log.New(os.Stdout, "", 0)
	} else {
		c.debug = log.New(io.Discard, "", 0)
	}
}

func (c *Chip8) GetDisplay() [32][64]bool {
	return c.display
}

func (c *Chip8) Debug() {
	c.debug.Println("pc:", c.pc)
	c.debug.Println("i:", c.i)
	c.debug.Println("vx:", c.register)
	c.debug.Println("stack:", c.stack)
	c.debug.Printf("sp: %d\n", c.sp)
}

func (c *Chip8) Fetch() uint16 {
	c.debug.Println("fetching: ", c.pc, "from ", len(c.memory))
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	c.pc += 2
	return opcode
}

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
		switch NN {
		case 0xE0: // clear Screen
			for i := 0; i < len(c.display); i++ {
				for j := 0; j < len(c.display[i]); j++ {
					c.display[i][j] = false
				}
			}
		case 0xEE: // return subroutine
			c.sp -= 1
			c.pc = c.stack[c.sp]
		}
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
		// Set
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
		}
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
		//display
		px := c.register[X] % 64
		py := c.register[Y] % 32
		c.register[0xF] = 0
		c.debug.Println("x,y: ", px, py)
		for row := 0; row < int(N); row++ {
			if int(py) >= len(c.display) { // reached bottom edge of screen
				break;
			}
			sbyte := c.memory[c.i+uint16(row)]
			c.debug.Printf("%s %x\n", "s: ", sbyte)
			px := c.register[X] % 64
			for bit := 0; bit < 8; bit++ {
				if int(px) >= len(c.display[0]) { // reached right edge of screen
					break;
				}
				spritePixel := (sbyte >> (7 - bit)) & 0x01
				// spritePixel := int(sbyte) & bit
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
	case 0xE000:
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
	case 0xF000:
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
			c.i = 0x50 + uint16(c.register[X]) * 5 // every font character is 5 bytes
		case 0x33: //binary-coded decimal conversion
			temp := c.register[X]
			i := 2
			for i >= 0 {
				digit := temp % 10
				c.memory[c.i+uint16(i)] = digit
				temp /= 10
				i--
			}
			// return fmt.Errorf("debug")
		case 0x55:
			for i := 0; i <= int(X); i += 1 {
				c.memory[c.i+uint16(i)] = c.register[i]
			}
		case 0x65:
			for i := 0; i <= int(X); i += 1 {
				c.register[i] = c.memory[c.i+uint16(i)]
			}
		}

	default:
		fmt.Printf("Invalid opcode %X\n", opcode)
	}

	if c.dt > 0 {
		c.dt -= 1
	}

	if c.st > 0 {
		c.st -= 1
	}
	return nil
}

func (c *Chip8) LoadProgram(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	if len(data) > len(c.memory)-0x200 { // program is loaded at 0x200
		return fmt.Errorf("program size %d bigger than available memory %d", len(data), len(c.memory)-0x200)
	}

	copy(c.memory[0x200:], data)
	return nil
}
