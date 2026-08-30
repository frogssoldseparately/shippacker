//go:build wasm

package mmrs

import (
	"fmt"
	"syscall/js"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

var fetch = js.Global().Get("fetch")

// Converts .mmrs to soundfont (if applicable) and sequence pair.
func RepackArchive(srcPath string, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, am *maps.Assets, bankId uint64) (uint16, error) {
	newFileCount := uint16(0)
	awaitable := fetch.Invoke(srcPath)
	ch := make(chan []js.Value)
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args
		return nil
	})
	cc := make(chan error)
	awaitable.Call("then", cb)
	go func() {
		results := <-ch
		rsp := results[0]
		awaitable = rsp.Call("bytes")
		go awaitable.Call("then", cb)
		results = <-ch
		fmt.Println(results[0].Length())
		rawArchiveBin := make([]byte, results[0].Length())
		for i := range rawArchiveBin {
			rawArchiveBin[i] = byte(results[0].Index(i).Int())
		}
		archive, err := sreader.OpenArchiveFromBytes(srcPath, &rawArchiveBin)
		if err != nil {
			fmt.Println("Failed at reader making")
			cc <- err
		} else {
			newFileCount, err = RepackArchiveFromZipReader(archive, lw, cw, am, bankId)
			cc <- err
		}
	}()
	outErr := <-cc
	return newFileCount, outErr
}
