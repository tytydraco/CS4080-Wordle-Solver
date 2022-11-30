package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const WORD_LEN = 5
const NUM_TRIES = 6

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

func GetWordFeedback(word string) []Feedback {
	letters := strings.Split(word, "")
	feedback := make([]Feedback, 5)

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

	return feedback
}

func RemoveInvalidWords(feedback []Feedback) int {
	// TODO: Go through all feedbacks (right place, wrong place, or not exists)
	//		 and remove words in the validWords that aren't possible answers.
	//		 Return the number of entries eliminated.

	return 0
}

func ChooseNextBestGuess() string {
	// TODO: Using some kind of heuristics, choose the next best word to pick

	return validWords[0]
}

func main() {
	validWords = GetValidWordList()
	fmt.Println(len(validWords))

	for i := 0; i < NUM_TRIES; i++ {
		fmt.Printf("Attempt %d/%d\n", i+1, NUM_TRIES)
		nextBestGuess := ChooseNextBestGuess()

		fmt.Printf("Best pick: %s\n", nextBestGuess)
		feedback := GetWordFeedback(nextBestGuess)
		// TODO: If all correct, we win!

		removed := RemoveInvalidWords(feedback)
		// TODO: Also remove our guess
		fmt.Printf("Eliminated %d words!\n", removed)
		fmt.Println()
	}
}
