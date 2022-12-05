package main

import (
	"fmt"
	"os"
	"strings"
)

const DEBUG = true
const WORD_LEN = 5
const NUM_TRIES = 6

// Empty struct indicating that in item exists in a set.
var exists = struct{}{}

// Calculate the frequencies that letters appear in each position in the list of valid words.
func GetLetterFrequencies() []map[rune]int {
	letterPositionFreqs := make([]map[rune]int, WORD_LEN)

	// Create a letter frequency map for each possible letter position.
	for i := 0; i < WORD_LEN; i++ {
		letterPositionFreqs[i] = make(map[rune]int)
	}

	// Sum up the number of letter occurances in each possible letter position.
	for _, word := range validWords {
		for position, letter := range word {
			letterPositionFreqs[position][letter]++
		}
	}

	return letterPositionFreqs
}

// Ask the user for which letters from the guess word were correct.
func GetWordFeedback(word string) []LetterCorrectness {
	letters := strings.Split(word, "")
	feedback := make([]LetterCorrectness, 5)

askWordFeedback:

	fmt.Printf("\nGive feedback for word: %s\n", word)
	fmt.Printf("--- ([c]orrect, [o]ut-of-order, [i]ncorrect, [q]uit, [r]eset) ---\n")

	// Get feedback for each letter.
	for i, letter := range letters {
	askLetterFeedback:
		var letterFeedbackStr string
		fmt.Printf("%s: ", letter)
		fmt.Scanf("%s\n", &letterFeedbackStr)

		// Allow the user to quit.
		if letterFeedbackStr == "q" {
			fmt.Println("Goodbye :)")
			os.Exit(0)
		}

		// Allow the user to reset their feedbacks.
		if letterFeedbackStr == "r" {
			goto askWordFeedback
		}

		letterCorrectness, isValidChar := letterCorrectnessMap[letterFeedbackStr]

		// If the user entered an invalid character, ask them again.
		for !isValidChar {
			fmt.Println("Bad feedback! Try again.")
			goto askLetterFeedback
		}

		feedback[i] = letterCorrectness
	}

	return feedback
}

// Given feedback from the user and the best guess we recommended, eliminate words that are definitely not the answer.
func RemoveInvalidWords(letterCorrectness []LetterCorrectness, bestGuess string) int {
	guessLetters := strings.Split(bestGuess, "")

	// Keep track of the letters that were not present in the word at all.
	incorrectLetters := make(map[string]struct{})
	for i, correctness := range letterCorrectness {
		if correctness == Incorrect {
			guessLetter := guessLetters[i]
			incorrectLetters[guessLetter] = exists
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
			if correctness == Correct && currentLetter != guessLetter {
				if DEBUG {
					fmt.Printf("[D] (1) removed '%s': '%s' does not match\n", validWord, currentLetter)
				}
				invalidWords[validWord] = exists
				removedWordsCount++
				break
			}

			// Checks if the current guess letter is supposed to be incorrectly positioned (but exists!), yet
			// matches the correct position in the current word.
			if correctness == OutOfOrder && currentLetter == guessLetter {
				if DEBUG {
					fmt.Printf("[D] (2) removed '%s': '%s' should not match\n", validWord, currentLetter)
				}

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
				if DEBUG {
					fmt.Printf("[D] (3) removed '%s': '%s' is incorrect and in order\n", validWord, currentLetter)
				}

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
			if DEBUG {
				fmt.Printf("[D] (4) removed '%s': not enough out-of-order chars (%d/%d)\n", validWord, outOfOrderLettersCount, len(outOfOrderChars))
			}

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

// Returns a map of words corresponding to its score in terms of how likely this word is to be the answer.
func GetWordScores(frequencies []map[rune]int) map[string]int {
	wordScores := make(map[string]int)

	// Sum up the frequency of letter occurances to determine a score.
	for _, word := range validWords {
		wordScore := 0
		for position, letter := range word {
			wordScore += frequencies[position][letter]
		}
		wordScores[word] = wordScore
	}

	return wordScores
}

// Calculate the highest score from possibleWords and returns the word with the highest score.
func GetNextBestGuess() string {
	letterFrequencies := GetLetterFrequencies()
	scores := GetWordScores(letterFrequencies)

	// Get the word with the highest score.
	maxScore := 0
	var bestWord string
	for word, score := range scores {
		if score > maxScore {
			maxScore = score
			bestWord = word
		}
	}

	return bestWord
}

// Return true if the user guessed the word.
func DidUserWin(feedback []LetterCorrectness) bool {
	for _, i := range feedback {
		if i != Correct {
			return false
		}
	}

	return true
}

func main() {
	UpdateValidWordsList()
	fmt.Printf("We have %d words to choose from...\n", len(validWords))

	// Try to guess the word in the limited number of tries.
	for i := 0; i < NUM_TRIES; i++ {
		fmt.Printf("Attempt %d/%d\n", i+1, NUM_TRIES)

		// Pick which word the user should guess.
		nextBestGuess := GetNextBestGuess()
		fmt.Printf("Best pick: %s\n", nextBestGuess)

		// Collect feedback on how we did.
		feedback := GetWordFeedback(nextBestGuess)

		// Check if the user won, and exit if they did.
		if DidUserWin(feedback) {
			fmt.Println("My work here is done :-)")
			fmt.Printf("Guessed the answer in %d/%d tries.\n", i+1, NUM_TRIES)
			break
		}

		// Tell the user how many words we were able to eliminate.
		removed := RemoveInvalidWords(feedback, nextBestGuess)
		fmt.Printf("Eliminated %d words!\n\n", removed)
	}
}
