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

### WASM

`wasm/main.go` compiles the `chip8` package to WebAssembly and exposes it to JavaScript as a global `Chip8` object (`load`, `step`, `tickTimers`, `needsRedraw`, `getDisplay`, `setKey`, `beeping`, `ips`). It holds no emulator logic itself — the browser side is expected to drive timing the same way `main.runLoop` does natively (run `ips()/60` `step()` calls per 60Hz tick, then `tickTimers()` once, then redraw if `needsRedraw()`).

```bash
GOOS=js GOARCH=wasm go build -o chip8.wasm ./wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .   # Go's JS glue for instantiating the module
```

## Feature Checklist

- [x] Implement functionality to display IBM emulator
- [x] Implement functionality to pass corax+ test
- [x] Add debug functionality
- [x] Flesh out API to control emulator for front-end/graphics
- [x] Add WASM compilation
- [ ] Stand up React/Typescript Front-End to consume WASM application
- [ ] Implement controls
- [ ] Debugging tools

## Known Issues

## TODO Roadmap

- [ ] Validate against the [Timendus chip8-test-suite](https://github.com/Timendus/chip8-test-suite) (flags test, quirks test) beyond corax+
- [ ] Add a `Quirks` struct to make behavior differences across CHIP-8 variants (shift semantics, `Fx55`/`Fx65` index increment, `BNNN` vs `BXNN`) configurable rather than hardcoded
- [ ] Add a graphical front-end (e.g. [Ebitengine](https://ebitengine.org), which also targets WASM) to replace the terminal `PrintDisplay` output
- [ ] Stand up the React/TypeScript front-end that consumes `wasm/`'s `Chip8` global.
