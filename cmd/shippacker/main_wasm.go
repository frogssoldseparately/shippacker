//go:build wasm

package main

import (
	"syscall/js"

	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
)

func main() {
	js.Global().Set("PackO2R", js.FuncOf(PackO2R))
	<-make(chan bool)
}

// js usage:
//
//	const bin = await PackO2R(...urls);
//	// do something with the zip binary
func PackO2R(this js.Value, args []js.Value) any {
	srcPaths := []string{}
	strArgs := args
	for _, arg := range strArgs {
		srcPaths = append(srcPaths, arg.String())
	}
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		go func() {
			goArr := shippacker.Pack(srcPaths)
			jsArr := js.Global().Get("Uint8Array").New(len(goArr))
			js.CopyBytesToJS(jsArr, goArr)
			resolve.Invoke(jsArr)
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
