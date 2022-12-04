package main

import (
	"bufio"
	"fmt"
	"math"
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
func WordScore(possibleWords []string, frequencies []map[byte]int) map[string]float64 {
	scores := make(map[string]float64)
	maxFreq := make([]int, WORD_LEN)

	// Get the max frequency in each position
	for i, v := range frequencies {
		for _, freq := range v {
			if freq > maxFreq[i] {
				maxFreq[i] = freq
			}
		}
	}

	// Calculate the score for each word by taking the difference of the
	// maximum frequency at each position and the frequency of the letter in the word at that position
	for _, word := range possibleWords {
		wordScore := float64(1)
		for i, v := range word {
			freqDiff := float64(frequencies[i][byte(v)] - maxFreq[i])
			wordScore += 1 + math.Pow(freqDiff, 2)
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
	for j, letter := range feedback {
		if letter == Nonexistent {
			nonexistentLetters[guessLetters[j]] = 0
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
			if nonexistent {
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
	maxScore := 0.0
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

func main() {
	validWords = GetValidWordList()
	unorderedWords = make(map[string]int)
	fmt.Println(len(validWords))

	for i := 0; i < NUM_TRIES; i++ {
		fmt.Printf("Attempt %d/%d\n", i+1, NUM_TRIES)
		nextBestGuess := ChooseNextBestGuess(validWords, LetterFrequency(validWords))

		fmt.Printf("Best pick: %s\n", nextBestGuess)
		feedback := GetWordFeedback(nextBestGuess)
		// TODO: If all correct, we win!

		removed := RemoveInvalidWords(feedback, nextBestGuess)
		// TODO: Also remove our guess
		fmt.Printf("Eliminated %d words!\n", removed)
		fmt.Println()
	}
}
