package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"unicode"
)

var (
	stdout = os.Stdout
	stderr = os.Stderr
	stdin  = os.Stdin
)

const (
	numLetters = 7
	minLength  = 4
)

func main() {
	fmt.Println("Enter 7 letters, the first letter is the required letter:")

	inputReader := bufio.NewReader(stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to get input: %s", err)
	}

	requiredLetter := rune(0)
	lettersMap := make(map[rune]struct{})
	for _, v := range input {
		if unicode.IsSpace(v) || v == ',' || v == ';' {
			continue
		}

		if requiredLetter == 0 {
			requiredLetter = v
		}
		lettersMap[v] = struct{}{}
	}
	if len(lettersMap) != numLetters {
		fmt.Fprintf(stderr, "Must specify exactly %d non-duplicate letters!\n", numLetters)
		os.Exit(127)
	}

	letters := make([]rune, 0, len(lettersMap))
	for letter := range lettersMap {
		letters = append(letters, letter)
	}
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
	fmt.Fprint(stdout, "Finding solutions for letters [")
	for i, letter := range letters {
		if i != 0 {
			fmt.Fprint(stdout, " ")
		}

		fmt.Fprint(stdout, string(letter))
		if letter == requiredLetter {
			fmt.Fprint(stdout, "*")
		}
	}
	fmt.Fprint(stdout, "]...\n\n")

	solutions, err := realMain(lettersMap, requiredLetter, words)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(stdout, "Found %d solutions:\n\n", len(solutions))

	sort.Slice(solutions, func(i, j int) bool {
		if len(solutions[j]) == len(solutions[i]) {
			return solutions[i] < solutions[j]
		}
		return len(solutions[j]) < len(solutions[i])
	})

	for _, solution := range solutions {
		fmt.Fprintf(stdout, "%s\n", solution)
	}
}

func realMain(lettersMap map[rune]struct{}, requiredLetter rune, dictionary []string) ([]string, error) {
	solutions := make([]string, 0, 32)

SCANNER:
	for _, word := range dictionary {
		// If the word is not long enough, it cannot be a solution
		if len(word) < minLength {
			continue
		}

		containsRequired := false
		for _, letter := range word {
			// If the word contains a letter that is not in our given set, it cannot
			// be a solution.
			if _, ok := lettersMap[letter]; !ok {
				continue SCANNER
			}

			if letter == requiredLetter {
				containsRequired = true
			}
		}

		// If we got this far, the world only contains allowed letters and contains
		// the required letter.
		if containsRequired {
			solutions = append(solutions, word)
		}
	}

	return solutions, nil
}
