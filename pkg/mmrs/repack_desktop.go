//go:build !wasm

package mmrs

import (
	"os"
	"path/filepath"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

// Converts .mmrs to soundfont (if applicable) and sequence pair.
func RepackArchive(musicSrcPath string, file os.DirEntry, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, am *maps.Assets, bankId uint64) (uint16, error) {
	archive, err := sreader.OpenArchive(filepath.Join(musicSrcPath, file.Name()))
	if err != nil {
		return 0, err
	}
	return RepackArchiveFromZipReader(archive, lw, cw, am, bankId)
}
