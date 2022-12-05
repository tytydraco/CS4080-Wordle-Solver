package main

import (
	"bufio"
	"os"
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

func UpdateValidWordsList() {
	validWords = GetValidWordList()
}
