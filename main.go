package main

import (
	"fmt"

	"github.com/chen-justin/chip-8/chip8"
)

func main() {

	c := chip8.Init()
	c.LoadProgram(".roms/3-corax+.ch8")
	// c.LoadProgram("./ibm.ch8")

	// c.Debug()
	// opcode := c.Fetch()
	// exec := c.Execute(opcode)
	// fmt.Println(exec)
	// PrintDisplay(c.GetDisplay())
	cycle := 0
	for {
		c.Debug()
		fmt.Println("cycle: ", cycle)
		opcode := c.Fetch()
		fmt.Printf("opcode: %x\n", opcode)
		c.Debug()
		e := c.Execute(opcode)
		if e != nil {
			break
		}
		cycle += 1
		PrintDisplay(c.GetDisplay())
		// time.Sleep(1000 / 60 * time.Millisecond) // Slow down output for visibility

	}

}

func PrintDisplay(display [32][64]bool) {

	for y := range display {
		for x := range display[y] {
			if display[y][x] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
