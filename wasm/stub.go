//go:build !(js && wasm)

// This file exists only so that `go build ./...`/`go vet ./...` from a
// native (non-wasm) target don't fail on a directory whose real file
// (main.go) is entirely excluded by its js/wasm build constraint. Build
// the actual WASM binary with:
//
//	GOOS=js GOARCH=wasm go build -o chip8.wasm ./wasm
package main

func main() {}
