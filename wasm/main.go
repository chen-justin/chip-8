//go:build js && wasm

// Command wasm compiles the chip8 package to WebAssembly and exposes it to
// JavaScript as a single global `Chip8` object. It holds no logic of its
// own beyond translating between JS values and the chip8.Chip8 API — the
// browser side drives timing (steps per tick, 60Hz TickTimers/redraw).
package main

import (
	"syscall/js"

	"github.com/chen-justin/chip-8/chip8"
)

var emulator chip8.Chip8

func main() {
	c := js.Global().Get("Object").New()
	c.Set("load", js.FuncOf(load))
	c.Set("step", js.FuncOf(step))
	c.Set("tickTimers", js.FuncOf(tickTimers))
	c.Set("needsRedraw", js.FuncOf(needsRedraw))
	c.Set("getDisplay", js.FuncOf(getDisplay))
	c.Set("setKey", js.FuncOf(setKey))
	c.Set("beeping", js.FuncOf(beeping))
	c.Set("ips", js.FuncOf(ips))
	js.Global().Set("Chip8", c)

	select {} // block forever: the registered funcs must stay alive
}

// load(romBytes: Uint8Array) -> {ok, error?}
// Resets the emulator and loads a fresh ROM. Call this before the first
// step() and again to restart.
func load(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errResult("load requires a Uint8Array ROM argument")
	}
	rom := args[0]
	data := make([]byte, rom.Get("length").Int())
	js.CopyBytesToGo(data, rom)

	emulator = chip8.Init()
	if err := emulator.LoadProgramBytes(data); err != nil {
		return errResult(err.Error())
	}
	return okResult()
}

// step() -> {ok, error?}
// Fetches and executes a single instruction. The caller is
// responsible for calling this IPS/60 times per 60Hz tick, then calling
// tickTimers() once.
func step(this js.Value, args []js.Value) any {
	if err := emulator.Step(); err != nil {
		return errResult(err.Error())
	}
	return okResult()
}

// tickTimers() -> undefined
func tickTimers(this js.Value, args []js.Value) any {
	emulator.TickTimers()
	return nil
}

// needsRedraw() -> boolean
func needsRedraw(this js.Value, args []js.Value) any {
	return emulator.NeedsRedraw()
}

// getDisplay() -> Uint8Array, 64*32 bytes, row-major, 0 or 1 per pixel
func getDisplay(this js.Value, args []js.Value) any {
	display := emulator.GetDisplay()
	pixels := make([]byte, len(display)*len(display[0]))
	for y := range display {
		for x := range display[y] {
			if display[y][x] {
				pixels[y*len(display[y])+x] = 1
			}
		}
	}

	out := js.Global().Get("Uint8Array").New(len(pixels))
	js.CopyBytesToJS(out, pixels)
	return out
}

// setKey(key: number, down: boolean) -> {ok, error?}
func setKey(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return errResult("setKey requires (key, down) arguments")
	}
	key := uint8(args[0].Int())
	down := args[1].Bool()
	if err := emulator.SetKey(key, down); err != nil {
		return errResult(err.Error())
	}
	return okResult()
}

// beeping() -> boolean
func beeping(this js.Value, args []js.Value) any {
	return emulator.Beeping()
}

// ips() -> number
func ips(this js.Value, args []js.Value) any {
	return emulator.IPS()
}

func okResult() map[string]any {
	return map[string]any{"ok": true}
}

func errResult(msg string) map[string]any {
	return map[string]any{"ok": false, "error": msg}
}
