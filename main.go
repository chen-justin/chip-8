package main

import (
	"flag"
	"fmt"
	"log"
	"time"

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

	if err := runLoop(&c, *debug); err != nil {
		fmt.Println(err)
	}
}

// runLoop drives the emulator at a fixed 60Hz: each tick runs a batch of
// IPS/60 instructions, then ticks the timers once and redraws at most once,
// decoupling both from however fast instructions actually execute.
func runLoop(c *chip8.Chip8, debug bool) error {
	const ticksPerSecond = 60
	instructionsPerTick := c.IPS() / ticksPerSecond

	ticker := time.NewTicker(time.Second / ticksPerSecond)
	defer ticker.Stop()

	cycle := 0
	for range ticker.C {
		for i := 0; i < instructionsPerTick; i++ {
			if debug {
				fmt.Println("cycle: ", cycle)
			}
			if err := c.Step(); err != nil {
				return fmt.Errorf("step error: %w", err)
			}
			cycle++
		}
		c.TickTimers()
		if c.NeedsRedraw() {
			PrintDisplay(c.GetDisplay())
		}
	}
	return nil
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
