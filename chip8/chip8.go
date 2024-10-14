package chip8

import (
	"fmt"
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
}

func Init() Chip8 {
	instance := Chip8{
		pc:     0x200,
		render: true,
		ips:    700,
	}
	buffer := 0x50
	for i := buffer; i < len(fontSet)+buffer; i++ {
		instance.memory[i] = fontSet[i-buffer]
	}
	return instance
}

func (c *Chip8) GetDisplay() [32][64]bool {
	return c.display
}

func (c *Chip8) Debug() {
	fmt.Println("pc:", c.pc)
	fmt.Println("i:", c.i)
	fmt.Println("vx:", c.register)
	fmt.Println("stack:", c.stack)
	fmt.Printf("sp: %d\n", c.sp)
}

func (c *Chip8) Fetch() uint16 {
	fmt.Println("fetching: ", c.pc, "from ", len(c.memory))
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

	fmt.Printf("%s %x\n", "nibble:", (n1))
	fmt.Printf("%s %d\n", "X: ", X)
	fmt.Printf("%s %d\n", "Y: ", Y)
	fmt.Printf("%s %d\n", "N: ", N)
	fmt.Printf("%s %x - %d\n", "NN: ", NN, NN)
	fmt.Printf("%s %x - %d\n", "NNN: ", NNN, NNN)
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
			// add
			c.register[X] = c.register[X] + c.register[Y]
		case 0x0005:
			// subtract
			c.register[X] = c.register[X] - c.register[Y]
		case 0x0007:
			// subtract
			c.register[X] = c.register[Y] - c.register[X]
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
		fmt.Println("x,y: ", px, py)
		for row := 0; row < int(N); row++ {
			if int(py) >= len(c.display) { // reached bottom edge of screen
				continue
			}
			sbyte := c.memory[c.i+uint16(row)]
			fmt.Printf("%s %x\n", "s: ", sbyte)
			px := c.register[X] % 64
			for bit := 0; bit < 8; bit++ {
				if int(px) >= len(c.display[px]) { // reached right edge of screen
					continue
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
		case 0x0A: // get key
			c.pc -= 2 // -1?
		case 0x29: //font character
			c.i = 0x50 + uint16(c.register[X]) //iffy
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
	file, fileErr := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	fStat, fStatErr := file.Stat()
	if fStatErr != nil {
		return fStatErr
	}
	fmt.Println("fileSize:", fStat.Size())
	if int64(len(c.memory)-512) < fStat.Size() { // program is loaded at 0x200
		return fmt.Errorf("program size bigger than memory")
	}

	buffer := make([]byte, fStat.Size())
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	fmt.Println("buffer", len(buffer))
	for i := 0; i < len(buffer); i++ {
		c.memory[i+512] = buffer[i]
	}
	fmt.Println("successfully loaded: ", fileName)
	return nil
}
