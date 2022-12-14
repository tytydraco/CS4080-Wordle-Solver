package main

type LetterCorrectness string

const (
	Correct    = "c"
	OutOfOrder = "o"
	Incorrect  = "i"
)

func (feedback LetterCorrectness) String() string {
	return string(feedback)
}

var (
	letterCorrectnessMap = map[string]LetterCorrectness{
		"c": Correct,
		"o": OutOfOrder,
		"i": Incorrect,
	}
)
