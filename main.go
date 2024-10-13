package main

import (
	"fmt"
)

const (
	width  = 64
	height = 32
)

func main() {
	// Create a 32x64 "display" initialized with all pixels off (false)
	display := make([][]bool, height)
	for i := range display {
		display[i] = make([]bool, width)
	}

	// Simulate turning on some pixels for demo purposes
	display[5][10] = true
	display[5][11] = true
	display[6][10] = true
	display[6][11] = true

	// Continuously print the display
	// for {
	// 	PrintDisplay(display)
	// 	// time.Sleep(500 * time.Millisecond) // Slow down output for visibility
	// }
	PrintDisplay(display)
}

func PrintDisplay(display [][]bool) {
	// Clear the console (simulated by printing new lines)
	fmt.Print("\033[H\033[2J")

	// Iterate over the display and print each row
	for _, row := range display {
		for _, pixel := range row {
			if pixel {
				fmt.Print("#") // Pixel on
			} else {
				fmt.Print(" ") // Pixel off
			}
		}
		fmt.Println()
	}
}
