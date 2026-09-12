package seqcat

import (
	"slices"
	"strings"
)

type listNode struct {
	Value string
	Next  *listNode
}

func (n *listNode) Append(str string) *listNode {
	newNode := &listNode{str, nil}
	n.Next = newNode
	return newNode
}

func ConvertToMMRSCategories(categories []string) *[]string {
	headNode := &listNode{"", nil}
	tailNode := headNode
	for _, cat := range categories {
		tailNode = tailNode.Append(cat)
	}
	intermediate := map[string]uint8{}
	headNode = headNode.Next
	for headNode != nil {
		curCat := headNode.Value
		if catArr, ok := ootCatConversion[curCat]; ok {
			for _, newCat := range catArr {
				tailNode = tailNode.Append(newCat)
			}
		} else if len(curCat) <= 3 {
			intermediate[curCat] = 0
		}
		headNode = headNode.Next
	}
	result := []string{}
	for key := range intermediate {
		result = append(result, key)
	}
	slices.Sort(result)
	return &result
}

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
	for _, cat := range mmCatFanfare {
		if slices.Contains(categories, cat) {
			return true
		}
	}
	return false
}
