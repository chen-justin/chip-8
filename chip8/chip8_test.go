package chip8

import "testing"

// newTestChip8 gives each test its own fresh, isolated instance.
func newTestChip8() Chip8 {
	return Init()
}

func TestLdVxByte(t *testing.T) {
	// 6XNN: VX = NN
	tests := []struct {
		name   string
		opcode uint16
		wantX  int
		wantNN uint8
	}{
		{"set V0 to 0x12", 0x6012, 0, 0x12},
		{"set VA to 0xFF", 0x6AFF, 0xA, 0xFF},
		{"set V0 to 0x00", 0x6000, 0, 0x00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChip8()
			if err := c.Execute(tt.opcode); err != nil {
				t.Fatalf("Execute(%#x) returned error: %v", tt.opcode, err)
			}
			if got := c.register[tt.wantX]; got != tt.wantNN {
				t.Errorf("register[%d] = %#x, want %#x", tt.wantX, got, tt.wantNN)
			}
		})
	}
}

func TestAddVxByte(t *testing.T) {
	// 7XNN: VX += NN. Per the spec this does NOT touch VF, even on overflow/wrap.
	tests := []struct {
		name    string
		startVX uint8
		nn      uint8
		wantVX  uint8
	}{
		{"simple add", 1, 2, 3},
		{"wraps on overflow", 0xFF, 2, 1}, // 0xFF + 2 = 0x101, truncates to 0x01
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChip8()
			c.register[0] = tt.startVX
			opcode := 0x7000 | uint16(tt.nn)
			if err := c.Execute(opcode); err != nil {
				t.Fatalf("Execute(%#x) returned error: %v", opcode, err)
			}
			if c.register[0] != tt.wantVX {
				t.Errorf("register[0] = %#x, want %#x", c.register[0], tt.wantVX)
			}
			if c.register[0xF] != 0 {
				t.Errorf("register[0xF] = %#x, want 0 (7XNN must not set VF)", c.register[0xF])
			}
		})
	}
}

// TestAddVxVyCarry documents the expected CHIP-8 behavior for 8XY4 (VX += VY),
// which must set VF to 1 on carry, 0 otherwise. This currently FAILS against
// the emulator's Execute implementation, which never touches VF for 8XY4 -
// that's a real bug to fix, not a mistake in the test.
func TestAddVxVyCarry(t *testing.T) {
	tests := []struct {
		name   string
		vx, vy uint8
		wantVX uint8
		wantVF uint8
	}{
		{"no carry", 10, 20, 30, 0},
		{"carry", 0xFF, 0x02, 0x01, 1}, // 0xFF + 0x02 = 0x101 -> VX=0x01, VF=1
		{"exact 255 no carry", 0x80, 0x7F, 0xFF, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChip8()
			c.register[0] = tt.vx
			c.register[1] = tt.vy
			if err := c.Execute(0x8014); err != nil { // VX=V0, VY=V1, N=4 (add)
				t.Fatalf("Execute(0x8014) returned error: %v", err)
			}
			if c.register[0] != tt.wantVX {
				t.Errorf("register[0] = %#x, want %#x", c.register[0], tt.wantVX)
			}
			if c.register[0xF] != tt.wantVF {
				t.Errorf("register[0xF] = %d, want %d", c.register[0xF], tt.wantVF)
			}
		})
	}
}

func TestSkipInstructions(t *testing.T) {
	tests := []struct {
		name     string
		opcode   uint16
		vx, vy   uint8
		wantSkip bool
	}{
		{"3XNN skips when equal", 0x30AA, 0xAA, 0, true},
		{"3XNN doesn't skip when not equal", 0x30AA, 0xAB, 0, false},
		{"4XNN skips when not equal", 0x40AA, 0xAB, 0, true},
		{"4XNN doesn't skip when equal", 0x40AA, 0xAA, 0, false},
		{"5XY0 skips when VX == VY", 0x5010, 5, 5, true},
		{"9XY0 skips when VX != VY", 0x9010, 5, 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChip8()
			startPC := c.pc
			c.register[0] = tt.vx
			c.register[1] = tt.vy
			if err := c.Execute(tt.opcode); err != nil {
				t.Fatalf("Execute(%#x) returned error: %v", tt.opcode, err)
			}
			wantPC := startPC
			if tt.wantSkip {
				wantPC += 2
			}
			if c.pc != wantPC {
				t.Errorf("pc = %d, want %d (skip=%v)", c.pc, wantPC, tt.wantSkip)
			}
		})
	}
}

func TestSetIndexRegister(t *testing.T) {
	// ANNN: I = NNN
	c := newTestChip8()
	if err := c.Execute(0xA123); err != nil {
		t.Fatalf("Execute(0xA123) returned error: %v", err)
	}
	if c.i != 0x123 {
		t.Errorf("i = %#x, want %#x", c.i, 0x123)
	}
}

func TestJumpAndCallReturn(t *testing.T) {
	c := newTestChip8()

	// 1NNN: unconditional jump
	if err := c.Execute(0x1300); err != nil {
		t.Fatalf("Execute(0x1300) returned error: %v", err)
	}
	if c.pc != 0x300 {
		t.Fatalf("after jump, pc = %#x, want %#x", c.pc, 0x300)
	}

	// 2NNN: call pushes return address and jumps
	if err := c.Execute(0x2400); err != nil {
		t.Fatalf("Execute(0x2400) returned error: %v", err)
	}
	if c.pc != 0x400 {
		t.Errorf("after call, pc = %#x, want %#x", c.pc, 0x400)
	}
	if c.sp != 1 {
		t.Fatalf("after call, sp = %d, want 1", c.sp)
	}
	if c.stack[0] != 0x300 {
		t.Errorf("stack[0] = %#x, want return address %#x", c.stack[0], 0x300)
	}

	// 00EE: return pops the stack back to the call site
	if err := c.Execute(0x00EE); err != nil {
		t.Fatalf("Execute(0x00EE) returned error: %v", err)
	}
	if c.pc != 0x300 {
		t.Errorf("after return, pc = %#x, want %#x", c.pc, 0x300)
	}
	if c.sp != 0 {
		t.Errorf("after return, sp = %d, want 0", c.sp)
	}
}

func TestClearScreen(t *testing.T) {
	c := newTestChip8()
	c.display[0][0] = true
	c.display[31][63] = true

	if err := c.Execute(0x00E0); err != nil {
		t.Fatalf("Execute(0x00E0) returned error: %v", err)
	}

	for y := range c.display {
		for x := range c.display[y] {
			if c.display[y][x] {
				t.Fatalf("display[%d][%d] = true, want false after clear", y, x)
			}
		}
	}
}

func TestGetKeyBlocksUntilPressed(t *testing.T) {
	c := newTestChip8()
	startPC := c.pc

	// No key pressed: Fx0A must rewind pc so the same instruction re-fetches
	// next cycle, and VX must stay untouched.
	if err := c.Execute(0xF00A); err != nil {
		t.Fatalf("Execute(0xF00A) returned error: %v", err)
	}
	if c.pc != startPC-2 {
		t.Errorf("pc = %d, want %d (should rewind when no key is pressed)", c.pc, startPC-2)
	}
	if c.register[0] != 0 {
		t.Errorf("register[0] = %#x, want 0 (VX should be untouched while waiting)", c.register[0])
	}

	// Restore pc to simulate the next fetch cycle re-running the instruction.
	c.pc = startPC

	// Now press key 0x7: VX should latch the key and execution should
	// proceed normally (no rewind).
	c.key[0x7] = true
	if err := c.Execute(0xF00A); err != nil {
		t.Fatalf("Execute(0xF00A) returned error: %v", err)
	}
	if c.register[0] != 0x7 {
		t.Errorf("register[0] = %#x, want %#x", c.register[0], 0x7)
	}
	if c.pc != startPC {
		t.Errorf("pc = %d, want %d (should not rewind once a key is pressed)", c.pc, startPC)
	}
}

func TestGetKeyPicksLowestPressedKey(t *testing.T) {
	c := newTestChip8()
	c.key[0x3] = true
	c.key[0x9] = true

	if err := c.Execute(0xF10A); err != nil { // target VX = V1
		t.Fatalf("Execute(0xF10A) returned error: %v", err)
	}
	if c.register[1] != 0x3 {
		t.Errorf("register[1] = %#x, want %#x (lowest pressed key)", c.register[1], 0x3)
	}
}

func TestUnknownSubOpcodeReturnsError(t *testing.T) {
	tests := []struct {
		name   string
		opcode uint16
	}{
		{"unknown 0x0NNN", 0x00FF},
		{"unknown 0x8XY?", 0x8009}, // nibble 9 isn't a defined 8XY? operation
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChip8()
			startPC := c.pc

			if err := c.Execute(tt.opcode); err == nil {
				t.Errorf("Execute(%#x) returned nil error, want non-nil", tt.opcode)
			}
			if c.pc != startPC {
				t.Errorf("pc = %d, want unchanged %d after an unknown opcode", c.pc, startPC)
			}
		})
	}
}

func TestReturnWithEmptyStackErrors(t *testing.T) {
	c := newTestChip8() // fresh instance: sp == 0, nothing has been called

	if err := c.Execute(0x00EE); err == nil {
		t.Error("Execute(0x00EE) with an empty call stack returned nil error, want non-nil")
	}
	if c.sp != 0 {
		t.Errorf("sp = %d, want unchanged 0 after a failed RET", c.sp)
	}
}

func TestFetch(t *testing.T) {
	t.Run("valid pc decodes and advances", func(t *testing.T) {
		c := newTestChip8() // pc starts at 0x200
		c.memory[0x200] = 0xAB
		c.memory[0x201] = 0xCD

		opcode, err := c.Fetch()
		if err != nil {
			t.Fatalf("Fetch() returned error: %v", err)
		}
		if opcode != 0xABCD {
			t.Errorf("opcode = %#x, want %#x", opcode, 0xABCD)
		}
		if c.pc != 0x202 {
			t.Errorf("pc = %#x, want %#x", c.pc, 0x202)
		}
	})

	t.Run("pc past end of memory errors without advancing", func(t *testing.T) {
		c := newTestChip8()
		c.pc = uint16(len(c.memory) - 1) // only one byte left: not enough for an opcode

		opcode, err := c.Fetch()
		if err == nil {
			t.Error("Fetch() returned nil error, want non-nil")
		}
		if opcode != 0 {
			t.Errorf("opcode = %#x, want 0", opcode)
		}
		if c.pc != uint16(len(c.memory)-1) {
			t.Errorf("pc = %#x, want unchanged %#x", c.pc, len(c.memory)-1)
		}
	})
}
