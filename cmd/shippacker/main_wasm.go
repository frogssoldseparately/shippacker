//go:build wasm

package main

import (
	"archive/zip"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/ootrs"
	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
	"github.com/frogssoldseparately/simpleseek/sreader"
)

func main() {
	js.Global().Set("PackO2R", js.FuncOf(PackO2R))
	js.Global().Set("AddOoTO2R", js.FuncOf(AddOoTO2R))
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

func AddOoTO2R(this js.Value, args []js.Value) any {
	buf := make([]byte, args[0].Length())
	js.CopyBytesToGo(buf, args[0])
	archive, err := sreader.OpenArchiveFromBytes("oot.o2r", &buf)
	if err != nil {
		fmt.Printf("Could not open archive because %s\n", err)
	}
	soundfontEntries := []*zip.File{}
	for _, entry := range archive.GetFiles() {
		if strings.Contains(entry.Name, "audio/fonts/") {
			soundfontEntries = append(soundfontEntries, entry)
		}
	}
	if err := ootrs.PrepareOotSoundfonts(&soundfontEntries); err == nil {
		globals.HasOotO2r = true
	} else {
		fmt.Printf("Could not unpack oot.o2r because %s\n", err)
	}
	return nil
}
