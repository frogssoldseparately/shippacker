//go:build !wasm

package ootrs

import (
	"os"
	"path/filepath"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func RepackArchive(musicSrcPath string, file os.DirEntry, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, am *maps.AssetMap, tm *maps.TranslationMap, bankId uint64) (uint16, error) {
	archive, err := sreader.OpenArchive(filepath.Join(musicSrcPath, file.Name()))
	if err != nil {
		return 0, err
	}
	return RepackArchiveFromZipReader(archive, lw, cw, am, tm, bankId)
}
