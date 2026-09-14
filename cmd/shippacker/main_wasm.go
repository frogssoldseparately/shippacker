//go:build wasm

package main

import (
	"fmt"
	"slices"
	"syscall/js"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/sample"
	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
	"github.com/frogssoldseparately/shippacker/pkg/soundfont"
	"github.com/frogssoldseparately/simpleseek/sreader"
)

func main() {
	js.Global().Set("PackO2R", js.FuncOf(PackO2R))
	js.Global().Set("ImportO2R", js.FuncOf(ImportO2R))
	js.Global().Set("SetO2RVersion", js.FuncOf(SetO2RVersion))
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
			if goArr := shippacker.Pack(srcPaths); goArr != nil {
				jsArr := js.Global().Get("Uint8Array").New(len(*goArr))
				js.CopyBytesToJS(jsArr, *goArr)
				resolve.Invoke(jsArr)
			} else {
				// return empty array
				jsArr := js.Global().Get("Uint8Array").New(0)
				resolve.Invoke(jsArr)
			}
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// js usage:
//
//	ImportO2R(o2rBinary);
//	const bin = await PackO2R(...urls);
//	// do something with the zip binary
func ImportO2R(this js.Value, args []js.Value) any {
	buf := make([]byte, args[0].Length())
	js.CopyBytesToGo(buf, args[0])
	archive, err := sreader.OpenArchiveFromBytes("import.o2r", &buf)
	if err != nil {
		fmt.Printf("Could not open archive because %s\n", err)
		globals.HasImportedO2R = false
		return nil
	}
	if err := soundfont.RegisterSoundfonts(archive); err != nil {
		fmt.Println(err)
		globals.HasImportedO2R = false
		return nil
	}
	sample.RegisterSamples(archive)
	globals.HasImportedO2R = true
	return nil
}

// js usage:
//
//	SetO2RVersion("version_string");
//	const bin = await PackO2R(...urls);
//	// do something with the zip binary
func SetO2RVersion(this js.Value, args []js.Value) any {
	v := args[0].String()
	if slices.Contains(globals.SupportedVersions, v) {
		globals.RomVersion = v
		if err := globals.SetupByVersion(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("%s is not a supported version\n", v)
	}
	return nil
}
