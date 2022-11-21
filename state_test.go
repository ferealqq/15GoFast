package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test State.newSwap
func TestNewSwap(t *testing.T){
	state := NewState()

	new := state.newSwap(len(state.board)-1,1)

	assert.Equal(t,new.board[1],0)
	assert.Equal(t,state.board[1],2)
}

func TestIsSolvable(t *testing.T){
	board := []int{
		12, 1, 10, 2,
		7, 11, 4, 14,
		5, 0, 9, 15, // Value 0 is used for empty space
		8, 13, 6, 3,	
	}
	assert.Equal(t, isSolvable(board),true)
	board = []int{
		3, 9, 1, 15,
		14, 11, 4, 6,
		13, 0, 10, 12,
		2, 7, 8, 5,
	}
	assert.Equal(t, isSolvable(board),false)

	board = []int{1,2,4,8,9,5,10,3,7,14,6,12,13,0,11,15}
	assert.Equal(t, isSolvable(board),true)
}

func TestGenerateState(t *testing.T){
	state,e := GenerateState(80)
	assert.Nil(t,e);

	fmt.Println(state.board);

	assert.True(t, isSolvable(state.board))

	// state,e = GenerateState(80)
	// assert.Nil(t,e)
	// assert.True(t, isSolvable(state.board))

	// state,e = GenerateState(38)
	// assert.Nil(t,e)
	// assert.True(t, isSolvable(state.board))

	// state,e = GenerateState(13)
	// assert.Nil(t,e)
	// assert.True(t, isSolvable(state.board))

	// state,e = GenerateState(69)
	// assert.Nil(t,e)
	// assert.True(t, isSolvable(state.board))
}

// get inversion count
func getInv(arr []int) int {
	total := BOARD_ROW_SIZE*BOARD_ROW_SIZE;
	inv := 0;
	for i := 0; i < total-1; i++ {
		for j := i + 1; j < total; j++ {
			if arr[i] > arr[j] {
				inv++;
			}
		}
	}
	return inv;
}

func isSolvable(board []int) bool {
	// Count inversions in given board
	invCount := getInv(board);
	// If grid is odd, return true if inversion
	// count is even.
	if BOARD_ROW_SIZE % 2 == 1 {
		return !(invCount % 2 == 1);
	}else{
		pos := GetElementIndex(board,0);
		fmt.Printf(" pos %d is even %d \n",pos, pos % 2)
		fmt.Printf(" inv count %d \n", invCount)
		if pos % 2 == 0 {
			return !(invCount % 2 == 0);
		}else{
			return invCount % 2 == 0;
		}
	}
}