package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInversionDistance(t *testing.T) {
	board := []int{1, 2, 4, 8, 9, 5, 10, 3, 7, 14, 6, 12, 13, 0, 11, 15}
	// state := NewSearch(NewState())
	inv1 := invertDistance(board)
	// https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm
	// TODO Not sure that the math behind this is correct, might have to check later
	assert.Equal(t, inv1, 32)
}

func TestHorizontal(t *testing.T) {
	board := []int{
		1, 2, 4, 8,
		9, 5, 10, 3,
		7, 14, 6, 12,
		13, 0, 11, 15,
	}

	trans := make([]int, len(board))
	copy(trans, board)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			trans[i*4+j] = board[i+j*4]
		}
	}

	to := []int{
		1, 9, 7, 13,
		2, 5, 14, 0,
		4, 10, 6, 11,
		8, 3, 12, 15,
	}

	assert.Equal(t, trans, to)
}

func TestSearchRandom(t *testing.T) {
	st, _ := GenerateState(7)
	board := st.board
	state := NewState()
	state.board = board
	fmt.Println("before state board")
	fmt.Println(state.board)
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	fmt.Println("after state board")
	fmt.Println(state.board)
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchEasy(t *testing.T) {
	board := []int{1, 2, 0, 4, 5, 6, 3, 8, 9, 10, 7, 11, 13, 14, 15, 12}
	state := NewState()
	state.board = board
	fmt.Println("before state board")
	fmt.Println(state.board)
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	fmt.Println("after state board")
	fmt.Println(state.board)
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchBug(t *testing.T) {
	board := []int{1, 2, 3, 4, 5, 6, 7, 8, 0, 9, 10, 12, 13, 14, 11, 15}	
	state := NewState()
	state.board = board
	fmt.Println("before state board")
	fmt.Println(state.board)
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	fmt.Println("after state board")
	fmt.Println(state.board)
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}