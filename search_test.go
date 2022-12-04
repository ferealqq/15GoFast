package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInversionDistance(t *testing.T) {
	board := startingPoint(4)
	for i, j := 0, len(board)-1; i < j; i, j = i+1, j-1 {
		board[i], board[j] = board[j], board[i]
	}
	inv1 := invertDistance(board)
	// https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm
	assert.Equal(t, inv1, t_cell(70))
}

func TestHorizontal(t *testing.T) {
	board := []t_cell{
		1, 2, 4, 8,
		9, 5, 10, 3,
		7, 14, 6, 12,
		13, 0, 11, 15,
	}

	trans := make([]t_cell, len(board))
	copy(trans, board)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			trans[i*4+j] = board[i+j*4]
		}
	}

	to := []t_cell{
		1, 9, 7, 13,
		2, 5, 14, 0,
		4, 10, 6, 11,
		8, 3, 12, 15,
	}

	assert.Equal(t, trans, to)
}

func TestSearchRandomFastEasy(t *testing.T) {
	st, _ := GenerateState(20)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchRandomFast40(t *testing.T) {
	st, _ := GenerateState(40)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestOptm(t *testing.T) {
	board := []t_cell{10, 9, 0, 4, 13, 11, 2, 8, 6, 3, 7, 12, 5, 1, 14, 15}
	fmt.Println(board)
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchRandomFastMedium(t *testing.T) {
	st, _ := GenerateState(50)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchRandomFastHard(t *testing.T) {
	st, _ := GenerateState(70)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

// named slow so that grouping test runs are easier
func TestSearchRandomSlowHard(t *testing.T) {
	t.Skip()
	st, _ := GenerateState(100)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchRandomSlowHarder(t *testing.T) {
	t.Skip()
	st, _ := GenerateState(150)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchRandomVeryHard(t *testing.T) {
	t.Skip()
	st, _ := GenerateState(300)
	board := st.board
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}


func TestSearchEasy(t *testing.T) {
	board := []t_cell{1, 2, 0, 4, 5, 6, 3, 8, 9, 10, 7, 11, 13, 14, 15, 12}
	state := NewState()
	state.board = board
	fmt.Println(state.board)
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	fmt.Println(state.board)
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))

	node.printMoves()
}

func TestSearchStuck(t *testing.T){
	board := []t_cell{5, 1, 2, 3, 9, 6, 8, 4, 13, 10, 7, 12, 14, 11, 15, 0}
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}

func TestSearchBugs(t *testing.T) {
	board := []t_cell{1, 2, 3, 4, 5, 6, 7, 8, 0, 9, 10, 12, 13, 14, 11, 15}
	state := NewState()
	state.board = board
	srh := NewSearch(state)
	node := srh.IDAStar(5)
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
	node.printMoves()
	board = []t_cell{5, 1, 2, 3, 9, 6, 8, 4, 13, 10, 7, 12, 14, 11, 15, 0}
	state = NewState()
	state.board = board
	srh = NewSearch(state)
	node = srh.IDAStar(5)
	node.printMoves()
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))

	board = []t_cell{5, 1, 2, 4, 9, 6, 3, 8, 13, 10, 7, 11, 0, 14, 15, 12}
	state = NewState()
	state.board = board
	srh = NewSearch(state)
	node = srh.IDAStar(5)
	assert.NotNil(t, node)
	assert.Equal(t, node.state.board, startingPoint(4))
}
