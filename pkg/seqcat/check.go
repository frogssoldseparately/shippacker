package seqcat

import (
	"slices"
	"strings"
)

func GetCategoriesFromString(str string) *[]string {
	categories := strings.FieldsFunc(str, func(r rune) bool { return r == ',' || r == '-' })
	return &categories
}

func IsFanfareString(catStr string) bool {
	switch catStr {
	case "fanfare":
		return true
	case "bgm":
		return false
	default:
		return HasFanfareCategories(*GetCategoriesFromString(catStr))
	}
}

func HasFanfareCategories(categories []string) bool {
	for _, cat := range fanfareCategories {
		if slices.Contains(categories, cat) {
			return true
		}
	}
	return false
}
