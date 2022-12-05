package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const WORD_LEN = 5
const NUM_TRIES = 6

// Empty struct indicating that in item exists in a set.
var exists = struct{}{}

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

// LetterFrequency calculates the frequency of each letter in each word of the possibleWords and
// returns a slice of maps, where each map is a map of the letter to its frequency in that position.
func LetterFrequency(possibleWords []string) []map[byte]int {
	freq := make([]map[byte]int, WORD_LEN)

	for i := 0; i < WORD_LEN; i++ {
		freq[i] = make(map[byte]int)
	}

	for _, word := range possibleWords {
		for i, v := range word {
			freq[i][byte(v)]++
		}
	}

	return freq
}

// WordScore calculates a score for each word in possibleWords based on the frequencies and
// returns a map of the word to its score.
func WordScore(possibleWords []string, frequencies []map[byte]int) map[string]int {
	scores := make(map[string]int)

	// Calculate the score for each word by summing the frequencies at each position
	for _, word := range possibleWords {
		wordScore := 0
		for i, v := range word {
			wordScore += frequencies[i][byte(v)]
		}
		scores[word] = wordScore
	}

	return scores
}

func GetWordFeedback(word string) []LetterCorrectness {
	letters := strings.Split(word, "")
	feedback := make([]LetterCorrectness, 5)

	for i, v := range letters {
	letterFeedback:
		var letterFeedbackStr string
		fmt.Printf("%s: ", v)
		fmt.Scanf("%s\n", &letterFeedbackStr)

		if letterFeedbackStr == "q" {
			fmt.Println("Goodbye :)")
			os.Exit(0)
		}

		letterFeedback, isPresent := letterCorrectnessMap[letterFeedbackStr]

		// Todo: make this better
		for !isPresent {
			fmt.Println("Bad feedback! Try again.")
			goto letterFeedback
		}

		feedback[i] = letterFeedback
	}

	return feedback
}

func RemoveInvalidWords(letterCorrectness []LetterCorrectness, bestGuess string) int {
	guessLetters := strings.Split(bestGuess, "")
	incorrectLetters := make(map[string]int)

	// Keep track of the letters that were not present in the word at all.
	for i, correctness := range letterCorrectness {
		if correctness == Incorrect {
			guessLetter := guessLetters[i]
			incorrectLetters[guessLetter] = 0
		}
	}

	// TODO(tytydraco): Refactor this bit
	//This loop checks to see if letters have been incorrectly added to the nonexsistenLetters map. If it has, it removes them from the map
	/*
		Example:
			Sores
			The first s is correct, the second s is nonexsistent
			Without this loop, s would be added to the nonexsistentLetters list and s would be considered to not exsist at all in the word
				This causes a problem, s IS in the word, it's just only present in the first space and not anywhere else
			This loop iterates through the word twice and checks if the letter occurs twice, and was given different ratings in both occurances
				If it does, and one of those instances it was rated as Nonexsistent, it is removed from the nonexsistent list as it does exsist, just not in that spot
	*/
	for i, correctness := range letterCorrectness {
		for j, checkLetter := range letterCorrectness {
			if guessLetters[i] == guessLetters[j] && ((correctness == Incorrect && checkLetter != Incorrect) || (checkLetter == Incorrect && correctness != Incorrect)) {
				delete(incorrectLetters, guessLetters[i])
			}
		}
	}

	// Checks which words from the possible picks list no longer work.
	invalidWords := make(map[string]struct{})
	outOfOrderChars := make(map[string]struct{})
	removedWordsCount := 0

	for _, validWord := range validWords {
		outOfOrderLettersCount := 0
		validWordLetters := strings.Split(validWord, "")
		for i, correctness := range letterCorrectness {
			currentLetter := validWordLetters[i]
			guessLetter := guessLetters[i]

			// Checks if the correct guess letter does not match the position in the current word.
			letterIsCorrect := correctness == Correct
			if letterIsCorrect && currentLetter != guessLetter {
				invalidWords[validWord] = exists
				removedWordsCount++
				break
			}

			// Checks if the current guess letter is supposed to be incorrectly positioned (but exists!), yet
			// matches the correct position in the current word.
			if correctness == OutOfOrder && currentLetter == guessLetter {
				invalidWords[validWord] = exists
				removedWordsCount++

				// Add this letter to the set of letters that are out of order.
				outOfOrderChars[currentLetter] = exists

				break
			}

			// TODO(tytydraco): Refactor this bit.
			// Checks if the current letter is not supposed to be in the word.
			_, letterIsIncorrect := incorrectLetters[currentLetter]
			_, letterIsOutOfOrder := outOfOrderChars[currentLetter]
			if letterIsIncorrect && !letterIsOutOfOrder {
				invalidWords[validWord] = exists
				removedWordsCount++
				break
			}

			// Checks if a nonexistent character is in the word
			if letterIsOutOfOrder {
				outOfOrderLettersCount++
			}
		}

		if len(outOfOrderChars) != 0 && outOfOrderLettersCount < len(outOfOrderChars) {
			invalidWords[validWord] = exists
			removedWordsCount++
		}
	}

	// Update the list of possible valid word picks.
	var newValidWords []string
	for _, v := range validWords {
		_, isInvalid := invalidWords[v]

		if !isInvalid {
			newValidWords = append(newValidWords, v)
		}
	}

	validWords = newValidWords
	return removedWordsCount
}

// ChooseNextBestGuess calculates the highest score from possibleWords and returns
// the word with the highest score.
func ChooseNextBestGuess(possibleWords []string, frequencies []map[byte]int) string {
	maxScore := 0
	bestWord := "words"
	scores := WordScore(possibleWords, frequencies)

	// Get the word with the highest score
	for word, score := range scores {
		if score > maxScore {
			maxScore = score
			bestWord = word
		}
	}

	return bestWord
}

func allCorrect(feedback []LetterCorrectness) bool {
	for _, i := range feedback {
		if i != Correct {
			return false
		}
	}

	return true
}

func main() {
	validWords = GetValidWordList()
	fmt.Println(len(validWords))

	for i := 0; i < NUM_TRIES; i++ {
		fmt.Printf("Attempt %d/%d\n", i+1, NUM_TRIES)
		nextBestGuess := ChooseNextBestGuess(validWords, LetterFrequency(validWords))

		fmt.Printf("Best pick: %s\n", nextBestGuess)
		feedback := GetWordFeedback(nextBestGuess)

		if allCorrect(feedback) {
			fmt.Println("My work here is done :-)")
			break
		}

		removed := RemoveInvalidWords(feedback, nextBestGuess)
		// TODO: Also remove our guess
		fmt.Printf("Eliminated %d words!\n", removed)
		fmt.Println()
	}
}
