//go:build !wasm

package ootrs

import (
	"os"
	"path/filepath"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func RepackArchive(musicSrcPath string, file os.DirEntry, zipWriter *swriter.SimpleZipWriter, gsm *maps.GameSampleMap) error {
	archive, err := sreader.OpenArchive(filepath.Join(musicSrcPath, file.Name()))
	if err != nil {
		return err
	}
	return RepackArchiveFromZipReader(archive, zipWriter, gsm)
}
