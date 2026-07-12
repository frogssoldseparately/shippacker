package seq

import (
	"fmt"
	"os"
	"strings"
)

type O2RMeta struct {
	DisplayName string
	BankIds     *[]byte
}

func NewMetaFromFile(path string) (*O2RMeta, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	str := string(raw)
	parts := strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")
	displayName := parts[0]
	if len(parts) == 1 {
		return nil, fmt.Errorf("Improperly formated .meta file.")
	}
	bankIds, err := parseBankHex(parts[1])
	for _, id := range *bankIds {
		if id > 0x28 {
			return nil, fmt.Errorf("Custom bank referenced but not through a .mmrs archive")
		}
	}
	if len(parts) == 2 || !strings.Contains(parts[2], "fanfare") {
		displayName += "_bgm"
	} else {
		displayName += "_fanfare"
	}
	return &O2RMeta{displayName, bankIds}, nil
}
