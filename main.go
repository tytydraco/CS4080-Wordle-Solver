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

// Contains a map of unordered words to choose from
var unorderedWords map[string]int

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

func GetWordFeedback(word string) []Feedback {
	letters := strings.Split(word, "")
	feedback := make([]Feedback, 5)

	for i, v := range letters {
		var letterFeedbackStr string
		fmt.Printf("%s: ", v)
		fmt.Scanf("%s\n", &letterFeedbackStr)

		letterFeedback, isPresent := feedbackMap[letterFeedbackStr]

		// Todo: make this better
		for !isPresent {
			fmt.Println("Bad feedback! Try again.")
			fmt.Printf("%s: ", v)
			fmt.Scanf("%s\n", &letterFeedbackStr)

			_, valid := feedbackMap[letterFeedbackStr]
			isPresent = valid
		}

		feedback[i] = letterFeedback
	}

	return feedback
}

func RemoveInvalidWords(feedback []Feedback, bestGuess string) int {
	// TODO: Go through all feedbacks (right place, wrong place, or not exists)
	//		 and remove words in the validWords that aren't possible answers.
	//		 Return the number of entries eliminated.
	guessLetters := strings.Split(bestGuess, "")
	nonexistentLetters := make(map[string]int)
	invalidWordsMap := make(map[string]int)

	var invalidWords int = 0

	//	Create a nonexistent letter map
	for j, letter := range feedback { //Adds non exsistent letters into nonexsistentLetters map
		if letter == Nonexistent {
			nonexistentLetters[guessLetters[j]] = 0
		}
	}

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
	for i, letter := range feedback {
		for j, checkLetter := range feedback {
			if guessLetters[i] == guessLetters[j] && ((letter == Nonexistent && checkLetter != Nonexistent) || (checkLetter == Nonexistent && letter != Nonexistent)) {
				delete(nonexistentLetters, guessLetters[i])
			}
		}
	}

	//	Goes through list of validWords, iterates through each character
	//	If any do not fit the feedback, it is added to the invalidWordsMap
	for _, v := range validWords {
		var unorderedWordCount int = 0
		validLetters := strings.Split(v, "")
		for j := range feedback {
			_, nonexistent := nonexistentLetters[validLetters[j]]
			_, letterIsPresent := unorderedWords[validLetters[j]]

			// Checks if the correct characters are in the correct place.
			if feedback[j] == Correct && validLetters[j] != guessLetters[j] {
				invalidWordsMap[v] = 0
				invalidWords++
				break
			}

			//	Checks if word has unordered characters
			if validLetters[j] == guessLetters[j] && feedback[j] == Unordered {
				invalidWordsMap[v] = 0
				invalidWords++
				if !letterIsPresent {
					unorderedWords[validLetters[j]] = 0
				}
				break
			}

			//	Checks if word contains nonexistent characters
			if nonexistent && !letterIsPresent {
				invalidWordsMap[v] = 0
				invalidWords++
				break
			}

			// Checks if a nonexistent character is in the word
			if letterIsPresent {
				unorderedWordCount++
			}
		}

		if unorderedWordCount < len(unorderedWords) && len(unorderedWords) != 0 {
			invalidWordsMap[v] = 0
			invalidWords++
		}
	}

	// Checks if word from validWords array exists in invalidWordsMap.
	// If not it does not add it to newValidWords array
	var newValidWords []string
	for _, v := range validWords {
		_, isPresent := invalidWordsMap[v]

		if !isPresent {
			newValidWords = append(newValidWords, v)
		}
	}

	validWords = newValidWords
	return invalidWords
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

func allCorrect(feedback []Feedback) bool {
	for _, i := range feedback {
		if i != Correct {
			return false
		}
	}

	return true
}

func main() {
	validWords = GetValidWordList()
	unorderedWords = make(map[string]int)
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
