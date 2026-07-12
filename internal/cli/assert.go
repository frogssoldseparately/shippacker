package cli

import (
	"os"
	"path/filepath"
)

func IsDirectory(param string, path string) error {
	if info, err := os.Stat(path); err != nil {
		return reportMissing(param, path)
	} else if !info.IsDir() {
		return reportNotDirectory(param, path)
	}
	return nil
}

func IsOfExtension(param string, path string, expectedExt string) error {
	if info, err := os.Stat(path); err != nil {
		return reportMissing(param, path)
	} else if info.IsDir() {
		return reportBadExtension(param, path, expectedExt, "Directory")
	} else if ext := filepath.Ext(info.Name()); ext != expectedExt {
		return reportBadExtension(param, path, expectedExt, ext)
	}
	return nil
}
