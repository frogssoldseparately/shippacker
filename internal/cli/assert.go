package cli

import "os"

func IsDirectory(param string, path string) error {
	if info, err := os.Stat(path); err != nil {
		return reportMissing(param, path)
	} else if !info.IsDir() {
		return reportNotDirectory(param, path)
	}
	return nil
}
