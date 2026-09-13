package iohelper

import (
	"fmt"
	"slices"
	"strings"
)

var validYesses = []string{"yes", "y"}

func PromptYesNo(msg string) bool {
	fmt.Printf("%s (y/n) ", msg)
	var response string
	for len(response) == 0 {
		fmt.Scan(&response)
		fmt.Scanln()
	}
	return slices.Contains(validYesses, strings.ToLower(response))
}

func PromptNumberedOption(msg string, opts []string, badMsg string) string {
	fmt.Println(msg)
	optCount := len(opts)
	for i, opt := range opts {
		fmt.Printf("\t%d) %s\n", i, opt)
	}
	response := -1
	for true {
		fmt.Print(":")
		fmt.Scan(&response)
		fmt.Scanln()
		if response >= 0 && response < optCount {
			break
		}
		fmt.Println(badMsg)
	}
	return opts[response]
}
