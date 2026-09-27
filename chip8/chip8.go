package chip8

import (
	"fmt"
	"io"
	"log"
	"os"
)

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
	loadFont(&instance.memory)
	return instance
}

func (c *Chip8) GetDisplay() [32][64]bool {
	return c.display
}

func (c *Chip8) Fetch() (uint16, error) {
	c.debug.Println("fetching: ", c.pc, "from ", len(c.memory))
	if int(c.pc)+1 >= len(c.memory) {
		return 0, fmt.Errorf("program counter out of bounds: %#04X", c.pc)
	}
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	c.pc += 2
	return opcode, nil
}

// Step fetches and executes a single instruction.
func (c *Chip8) Step() error {
	opcode, err := c.Fetch()
	if err != nil {
		return err
	}
	return c.Execute(opcode)
}

// TickTimers decrements the delay and sound timers by one, flooring at zero.
// Call this at a fixed 60Hz, independent of how fast instructions run.
func (c *Chip8) TickTimers() {
	if c.dt > 0 {
		c.dt--
	}
	if c.st > 0 {
		c.st--
	}
}

// NeedsRedraw reports whether the display has changed since the last call,
// clearing the flag as it does.
func (c *Chip8) NeedsRedraw() bool {
	needs := c.render
	c.render = false
	return needs
}

// IPS returns the configured instructions-per-second rate.
func (c *Chip8) IPS() int {
	return c.ips
}

func (c *Chip8) LoadProgram(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return c.LoadProgramBytes(data)
}

func (c *Chip8) LoadProgramBytes(data []byte) error {
	if len(data) > len(c.memory)-0x200 { // program is loaded at 0x200
		return fmt.Errorf("program size %d bigger than available memory %d", len(data), len(c.memory)-0x200)
	}

	copy(c.memory[0x200:], data)
	return nil
}

func (c *Chip8) SetKey(k uint8, down bool) error {
	if int(k) >= len(c.key) {
		return fmt.Errorf("key %#X out of range (must be 0x0-0xF)", k)
	}
	c.key[k] = down
	return nil
}

func (c *Chip8) Beeping() bool {
	return c.st > 0
}
