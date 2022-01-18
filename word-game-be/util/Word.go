package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var wordList = NewWordList()

type WordList struct {
	Words []string
	Index map[rune]int
}

func NewWordList() *WordList {
	words, err := os.Open("words.txt")

	if err != nil {
		panic(err)
	}

	wl := WordList{
		Words: make([]string, 0),
		Index: make(map[rune]int),
	}

	lastLetter := ' '
	ind := 0

	buffer := bufio.NewReader(words)
	var line string
	for {
		line, err = buffer.ReadString('\n')
		if err != nil || line == "" {
			break
		}
		letter := rune(line[0])
		ind++
		if letter != lastLetter {
			wl.Index[letter] = ind
		}
		lastLetter = letter
		wl.Words = append(wl.Words, strings.Trim(line, "\n"))

	}

	fmt.Printf("Init wordList with %d words\n", len(wl.Words))

	return &wl
}

func IsValidWord(word string) bool {
	word = strings.ToUpper(word)
	initial := rune(word[0])
	startIndex := wordList.Index[initial]
	for i := startIndex; i < len(wordList.Words); i++ {
		dictWord := wordList.Words[i]
		if rune(dictWord[0]) != initial {
			break
		}
		if dictWord == word {
			return true
		}
	}

	return false
}
