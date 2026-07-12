package cli

import (
	"fmt"
)

func reportMissing(param string, path string) error {
	return fmt.Errorf("(arg: %s) %s does not exist\n", param, path)
}

func reportNotDirectory(param string, path string) error {
	return fmt.Errorf("(arg: %s) %s is not a directory\n", param, path)
}

func reportBadExtension(param string, path string, expected string, actual string) error {
	return fmt.Errorf("(arg: %s) %s is a %s file instead of a %s file\n", param, path, actual, expected)
}
