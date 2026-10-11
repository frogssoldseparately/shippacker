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
