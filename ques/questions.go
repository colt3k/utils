package ques

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/colt3k/nglog/ng"
)

// Question used for questions that require a string response
func Question(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	selection, _ := reader.ReadString('\n')
	if strings.Index(selection, "\r") > -1 {
		selection = strings.TrimSuffix(selection, "\r\n")
	} else {
		selection = strings.TrimSuffix(selection, "\n")
	}

	return strings.TrimSpace(selection)
}

func QuestionOpts(question string, options []string) string {
optLoop:
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	selection, _ := reader.ReadString('\n')
	if strings.Index(selection, "\r") > -1 {
		selection = strings.TrimSuffix(selection, "\r\n")
	} else {
		selection = strings.TrimSuffix(selection, "\n")
	}

	selection = strings.TrimSpace(selection)
	var found bool
	for _, d := range options {
		if selection == d {
			found = true
			break
		}
	}
	if !found {
		log.Println("invalid option")
		goto optLoop
	}
	return strings.TrimSpace(selection)
}

func QuestionInt(question string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	selection, _ := reader.ReadString('\n')
	if strings.Index(selection, "\r") > -1 {
		selection = strings.TrimSuffix(selection, "\r\n")
	} else {
		selection = strings.TrimSuffix(selection, "\n")
	}

	regexp.MustCompile("")
	number, err := strconv.Atoi(strings.TrimSpace(selection))
	if err != nil {
		log.Logln(log.ERROR, err)
	}
	return number
}

// Confirm used for yes/no, true/false questions.
func Confirm(question string) bool {
	valid := []string{"t", "y", "yes", "true", "f", "n", "no", "false"}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	selection, _ := reader.ReadString('\n')
	if strings.Index(selection, "\r") > -1 {
		selection = strings.TrimSuffix(selection, "\r\n")
	} else {
		selection = strings.TrimSuffix(selection, "\n")
	}

	var validAns bool
	for _, d := range valid {
		if strings.EqualFold(strings.TrimSpace(selection), d) {
			validAns = true
		}
	}
	if !validAns {
		log.Logln(log.ERROR, "invalid input, please try again")
		return Confirm(question)
	}
	if strings.EqualFold(strings.TrimSpace(selection), "y") ||
		strings.EqualFold(strings.TrimSpace(selection), "true") ||
		strings.EqualFold(strings.TrimSpace(selection), "yes") ||
		strings.EqualFold(strings.TrimSpace(selection), "t") {
		return true
	}
	return false
}
