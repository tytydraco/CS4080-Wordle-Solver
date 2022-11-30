package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const WORD_LEN = 5

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
	feedback := new([WORD_LEN]string)

	fmt.Printf(" --- Feedback for %s --- \n", word)
	for i, v := range letters {
		var letterFeedback string
		fmt.Printf("%s: ", v)
		fmt.Scanf("%s", &letterFeedback)

		feedback[i] = letterFeedback
	}

	fmt.Println(feedback)
}

func main() {
	words := GetValidWordList()
	fmt.Println(words)

	GetWordFeedback(words[0])
}
