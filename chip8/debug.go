package chip8

import (
	"io"
	"log"
	"os"
)

// SetDebug turns verbose per-instruction tracing to stdout on or off.
func (c *Chip8) SetDebug(on bool) {
	if on {
		c.debug = log.New(os.Stdout, "", 0)
	} else {
		c.debug = log.New(io.Discard, "", 0)
	}
}

func (c *Chip8) Debug() {
	c.debug.Println("pc:", c.pc)
	c.debug.Println("i:", c.i)
	c.debug.Println("vx:", c.register)
	c.debug.Println("stack:", c.stack)
	c.debug.Printf("sp: %d\n", c.sp)
}
