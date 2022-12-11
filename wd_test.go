package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// walking distance is based on
// https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm

func TestWdReverse(t *testing.T) {
	board := startingPoint(BOARD_ROW_SIZE)
	for i, j := 0, len(board)-1; i < j; i, j = i+1, j-1 {
		board[i], board[j] = board[j], board[i]
	}

	wd := NewWD(int(BOARD_ROW_SIZE))
	assert.Equal(t, 70, wd.Calculate(board))
}

func TestWd(t *testing.T) {
	board := [16]t_cell{1, 2, 3, 0, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 4}

	wd := NewWD(int(BOARD_ROW_SIZE))
	assert.Equal(t, 15, wd.Calculate(board))
}
func TestWdMaybe(t *testing.T) {
	wd := NewWD(int(BOARD_ROW_SIZE))
	// these calculations are based on the older version of the walking distance
	// so these heuristic values could be wrong
	board := [16]t_cell{5, 1, 2, 4, 9, 6, 3, 8, 13, 10, 7, 11, 0, 14, 15, 12}
	assert.Equal(t, wd.Calculate(board), 9)
	board = [16]t_cell{1, 2, 0, 4, 5, 6, 3, 8, 9, 10, 7, 11, 13, 14, 15, 12}
	assert.Equal(t, wd.Calculate(board), 4)
	board = [16]t_cell{1, 2, 3, 4, 0, 5, 7, 8, 10, 6, 11, 12, 9, 13, 14, 15}
	assert.Equal(t, wd.Calculate(board), 7)
	board = [16]t_cell{1, 2, 3, 4, 5, 7, 11, 8, 9, 6, 14, 12, 13, 10, 15, 0}
	assert.Equal(t, wd.Calculate(board), 8)
	board = [16]t_cell{1, 2, 3, 4, 5, 6, 7, 8, 0, 9, 10, 12, 13, 14, 11, 15}
	assert.Equal(t, wd.Calculate(board), 4)
}

func TestWdGenerateTable(t *testing.T) {
	wd := NewWD(int(BOARD_ROW_SIZE))

	// based on https://computerpuzzle.net/english/15puzzle/wd.gif
	assert.Equal(t, len(wd.table), 24964)
	cc := 0
	for _, j := range wd.table {
		if j == 0 {
			cc++
		}
	}

	assert.Equal(t, cc, 1)
}
