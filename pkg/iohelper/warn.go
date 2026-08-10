package iohelper

import (
	"fmt"
)

const (
	ContinueRunning int = iota
	EarlyExit
	IgnoreOtherBanks
	HaltRunning
)

func promptUser(msg string) string {
	fmt.Printf("%s\n:", msg)
	var out string
	fmt.Scan(&out)
	fmt.Scanln() // discard newline?
	return out
}

func WarnPromptBanks() int {
	userInput := promptUser("WARNING: The limit of custom sound fonts has been reached (255). Adding any more can lead to audio bugs in 2ship.\n\nPlease choose one of the following:\n\tContinue anyway: type \"continue\"\n\tContinue and ignore future custom banks: type \"ignore\"\n\tGenerate O2R now: type \"finish\"\n\tCancel: type \"quit\"")
	switch userInput {
	case "continue":
		return ContinueRunning
	case "ignore":
		return IgnoreOtherBanks
	case "finish":
		return EarlyExit
	case "quit":
		fallthrough
	default:
		fmt.Println("HaltRunning")
		return HaltRunning
	}
}

func WarnPromptSongs() int {
	userInput := promptUser("WARNING: The limit of custom sequences has been reached (1923). Adding any more will lead to audio bugs in 2ship.\n\nIf you want to continue adding songs, type \"continue\", or \"no\" to stop")
	if userInput != "continue" {
		return EarlyExit
	}
	return ContinueRunning
}

func WarnUnstable() error {
	userInput := promptUser("WARNING: The selected version may generate o2rs not supported by the latest release of 2ship.\n\nIf you are working on a branch of 2ship that explicitly supports this feature, type \"continue\", or \"no\" to stop")
	if userInput != "continue" {
		return fmt.Errorf("Exiting. Please run again")
	}
	return nil
}
