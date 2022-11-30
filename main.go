package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const WORD_LEN = 5

type Feedback string

const (
	Correct     = "c"
	Unordered   = "u"
	Nonexistent = "n"
)

func (feedback Feedback) String() string {
	return string(feedback)
}

var (
	feedbackMap = map[string]Feedback{
		"c": Correct,
		"u": Unordered,
		"n": Nonexistent,
	}
)

// List of valid words to choose from
var validWords []string

func GetValidWordList() []string {
	readFile, _ := os.Open("words.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var fileLines []string
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

func GetWordFeedback(word string) {
	letters := strings.Split(word, "")
	feedback := new([WORD_LEN]Feedback)

	fmt.Printf(" --- Feedback for %s --- \n", word)
	for i, v := range letters {
		var letterFeedbackStr string
		fmt.Printf("%s: ", v)
		fmt.Scanf("%s", &letterFeedbackStr)

		letterFeedback, exists := feedbackMap[letterFeedbackStr]

		// Todo: make this better
		if !exists {
			fmt.Println("Bad feedback!")
			os.Exit(1)
		}

		feedback[i] = letterFeedback
	}

	fmt.Println(feedback)
}

func RemoveInvalidWords(feedback Feedback) {
	// TODO: Go through all feedbacks (right place, wrong place, or not exists)
	//		 and remove words in the validWords that aren't possible answers
}

func ChooseNextBestGuess() string {
	// TODO: Using some kind of heuristics, choose the next best word to pick

	return validWords[0]
}

func main() {
	validWords = GetValidWordList()
	fmt.Println(validWords)

	GetWordFeedback(validWords[0])
}
