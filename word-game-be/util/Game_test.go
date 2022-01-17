package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWordBoundsHoz(t *testing.T) {
	// Straight forward
	testHozLayout(t, "----ABC-", 4, 2, 4, 6)
	// Full line
	testHozLayout(t, "ABCD----", 4, 2, 0, 3)
	// Full line w/ other tiles
	testHozLayout(t, "ABCD-DEF", 4, 2, 0, 3)
	// Breaking over multiple lines
	testHozLayout(t, "--ABCD--", 4, 2, 2, 3)
	// Single letter
	testHozLayout(t, "--A-----", 4, 2, 2, 2)
}

func testHozLayout(t *testing.T, layout string, width int, height int, wordStart int, wordEnd int) {
	board := NewBoardFromLayout(layout, width, height)

	for i := wordStart; i <= wordEnd; i++ {
		start, end := GetWordBoundsHoz(board, i)
		assert.Equal(t, wordStart, start, "Start calculated incorrectly from target %d", i)
		assert.Equal(t, wordEnd, end, "End calculated incorrectly from target %d", i)
	}
	fmt.Println("---")
}

func testVertLayout(t *testing.T, layout string, width int, height int, wordStart int, wordEnd int) {
	board := NewBoardFromLayout(layout, width, height)

	for i := wordStart; i <= wordEnd; i += width {
		start, end := GetWordBoundsVert(board, i)
		assert.Equal(t, wordStart, start, "Start calculated incorrectly from target %d", i)
		assert.Equal(t, wordEnd, end, "End calculated incorrectly from target %d", i)
		fmt.Println("---")
	}
	fmt.Println("===")
}

func TestGetWordBoundsVert(t *testing.T) {
	testVertLayout(t, "-A---B---C---D--", 4, 4, 1, 13)
	testVertLayout(t, "A--B--C--D--", 3, 4, 0, 10)
}
