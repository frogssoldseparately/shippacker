//go:build wasm

package main

import (
	"syscall/js"

	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
)

var zipOut []byte

func main() {
	js.Global().Set("PackO2R", js.FuncOf(Pack))
	<-make(chan bool)
}

func Pack(this js.Value, args []js.Value) any {
	srcPaths := []string{}
	for _, arg := range args {
		srcPaths = append(srcPaths, arg.String())
	}
	zipOut = shippacker.Pack(srcPaths)
	return nil
}
