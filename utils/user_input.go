package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func Confirm(prompt string, defaultYes bool) bool {
	var yn string
	if defaultYes {
		yn = "Y/n"
	} else {
		yn = "y/N"
	}
	question := fmt.Sprintf("%s (%s)?", prompt, yn)
	text, err := GetInput(question)
	if err != nil {
		return false
	}

	text = strings.ToLower(text)
	if text == "" && defaultYes {
		return true
	}
	if strings.HasPrefix(text, "y") {
		return true
	}

	return false
}

func GetInput(prompt string) (string, error) {
	if prompt != "" {
		fmt.Println(prompt)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Trim the newline from the end
	return text[:len(text)-1], nil
}

func GetPasswordInput(prompt string) ([]byte, error) {
	if prompt != "" {
		fmt.Println(prompt)
	}

	fmt.Printf("> ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	fmt.Println("")

	return password, nil
}
