package seq

import (
	"fmt"
	"os"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/seqcat"
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
	suffix := "bgm"
	if len(parts) >= 3 {
		catString := strings.Trim(strings.ToLower(parts[2]), " \t")
		if catString == "fanfare" {
			suffix = "fanfare"
		} else if catString != "bgm" && len(catString) > 0 {
			categories := seqcat.GetCategoriesFromString(catString)
			// TODO: clean up this absolute mess
			if globals.UseNumericCategories {
				suffix = strings.Join(*categories, "-")
			} else if seqcat.HasFanfareCategories(*categories) {
				suffix = "fanfare"
			}
		}
	}
	displayName += "_" + suffix
	return &O2RMeta{displayName, bankIds}, nil
}
