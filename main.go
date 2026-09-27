package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chen-justin/chip-8/chip8"
)

func main() {
	romPath := flag.String("rom", ".roms/3-corax+.ch8", "path to a CHIP-8 ROM file to load")
	debug := flag.Bool("debug", false, "print verbose per-instruction debug output")
	flag.Parse()

	c := chip8.Init()
	c.SetDebug(*debug)
	if err := c.LoadProgram(*romPath); err != nil {
		log.Fatalf("failed to load ROM %q: %v", *romPath, err)
	}

	cycle := 0
	for {
		opcode, ferr := c.Fetch()
		if ferr != nil {
			fmt.Println("fetch error:", ferr)
			break
		}
		if *debug {
			fmt.Println("cycle: ", cycle)
			fmt.Printf("opcode: %x\n", opcode)
		}
		c.Debug()
		e := c.Execute(opcode)
		if e != nil {
			fmt.Println("execute error:", e)
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
