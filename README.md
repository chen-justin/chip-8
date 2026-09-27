# chip-8

A WIP project to implement a chip-8 interpreter/emulator in Golang.

Uses direction from this great guide: https://tobiasvl.github.io/blog/write-a-chip-8-emulator/#00e0-clear-screen

## Usage

Run with a ROM path via the `-rom` flag; it defaults to `.roms/3-corax+.ch8` if omitted.

```bash
go run . -rom .roms/3-corax+.ch8
go run . -h   # list flags
```

The [Timendus chip8-test-suite](https://github.com/Timendus/chip8-test-suite) ROMs live under `.roms/chip8-test-suite-4.2/bin/`, one file per stage (`1-chip8-logo.ch8` through `8-scrolling.ch8`):

```bash
go run . -rom .roms/chip8-test-suite-4.2/bin/4-flags.ch8
go run . -rom .roms/chip8-test-suite-4.2/bin/5-quirks.ch8
```

`6-keypad.ch8` needs real key input to progress past its first screen, which isn't wired up yet (see Roadmap item 6).

## Feature Checklist

- [x] Implement functionality to display IBM emulator
- [x] Implement functionality to pass corax+ test
- [ ] Add debug functionality
- [ ] Flesh out API to control emulator for front-end/graphics
- [ ] Add WASM compilation
- [ ] Stand up React/Typescript Front-End to consume WASM application
- [ ] Implement controls
- [ ] Debugging tools

## Known Issues

- [ ] Timers (`dt`/`st`) decrement once per instruction in `Execute`, instead of at a fixed 60Hz independent of CPU speed
- [ ] Unknown/invalid sub-opcodes (bad `0x0NNN`, `0x8XY?`) and stack underflow on `00EE` fail silently instead of returning an error
- [ ] `render` and `ips` fields on `Chip8` are unused
- [ ] Debug `fmt.Print` calls inside `Fetch`/`Execute` run unconditionally instead of being gated behind a debug flag

## Roadmap

1. [x] Add opcode tests (`chip8_test.go`) before refactoring, so later changes can be made safely
2. [ ] Fix the bugs listed above under Known Issues, confirming each fix against the tests
3. [ ] Validate against the [Timendus chip8-test-suite](https://github.com/Timendus/chip8-test-suite) (flags test, quirks test) beyond corax+
4. [ ] Add a `Quirks` struct to make behavior differences across CHIP-8 variants (shift semantics, `Fx55`/`Fx65` index increment, `BNNN` vs `BXNN`) configurable rather than hardcoded
5. [ ] Split opcode decoding from execution (e.g. per-family methods like `op8XY`, `opF`) to keep `Execute` from growing into one large switch
6. [ ] Define a stable front-end-facing API on `Chip8`: `Step()`, `TickTimers()`, `SetKey(k uint8, down bool)`, `Display()`, `Beeping()`
7. [ ] Build a real run loop: a 60Hz `time.Ticker` driving ~11 `Step()` calls per tick (700 IPS), then a draw call
8. [ ] Add a graphical front-end (e.g. [Ebitengine](https://ebitengine.org), which also targets WASM) to replace the terminal `PrintDisplay` output
