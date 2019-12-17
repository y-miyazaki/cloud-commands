package main

import (
	"fmt"
	"os/exec"
)

// execCommand execute exec.Command and output command.
func execCommand(command string, options ...string) (string, error) {
	out, err := exec.Command(command, options...).CombinedOutput()
	outputCommand := command + " "
	for _, s := range options {
		outputCommand = s + " "
	}
	fmt.Println(outputCommand)
	fmt.Println(string(out))
	return string(out), err
}

// execCommandStr execute exec.Command and output command.
func execCommandStr(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()
	fmt.Println(command)
	fmt.Println(string(out))
	return string(out), err
}
