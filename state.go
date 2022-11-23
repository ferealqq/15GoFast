package main

import (
	b64 "encoding/base64"
	"errors"
	"math/rand"
	"time"
)

const BOARD_ROW_SIZE = 4

type State struct {
	size int
	board []int
	complexity int
}

func NewState() *State{
	return &State{ 
		size: BOARD_ROW_SIZE, 
		board: startingPoint(BOARD_ROW_SIZE),
		complexity: 0,
	}
}

// Because all boards are not solvable we have to make a board that is solvalable 
func GenerateState(complexity int) (*State, error) {
	// https://www.geeksforgeeks.org/check-instance-15-puzzle-solvable/
	state := NewState()
	visited := []string{}
	olds := []*State{}

	rand.Seed(time.Now().UnixNano())
	for state.complexity < complexity {
		visited = append(visited, state.hash())
		sts := state.GetValidStates()
		filtered := []*State{}
		for i := range sts {
			for j := range visited {
				if visited[j] != sts[i].hash() {
					filtered = append(filtered, sts[i])
					break
				}
			}
		}
		if len(filtered) == 0 {
			return nil, errors.New("Too complex board was about to be created")
		}
		rn := rand.Intn(len(filtered))
		olds = append(olds, state)
		state = filtered[rn]
	}

	return state, nil
}

func startingPoint(size int) []int {
	res := make([]int, size*size)

	for i := 0; int(i) < (size*size)-1; i++ {
		res[i] = int(i+1);
	}

	return res;
}

// Get the first index of a element in a given array, returns -1 if not found
func GetElementIndex[T comparable](arr []T,element T) int {
	// get empty index 
	for i := 0; i < len(arr); i++ {
		if arr[i] == element {
			return i;
		}
	}
	return -1;
}

// Calculate all the states that the current state can be mutated to
func (state *State) GetValidStates() []*State {
	var moves = []*State{}

	// get emtpy index from the board
	emptyIndex := GetElementIndex(state.board, 0) 
	
	// the basic logic of puzzle 15 game is that you cannot move off the grid 16x16
	// and you have to move one step at the time at 16x16 grid. Which means that you have to move either +1 -1 +4 -4
	// for example you cannot move up if you are in the top row of the grid. 
	if (emptyIndex - state.size >= 0) {
		moves = append(moves, state.newSwap(emptyIndex, emptyIndex - state.size))
	}
	// Not on last line
	if (emptyIndex + state.size < len(state.board)) {
		moves = append(moves, state.newSwap(emptyIndex, emptyIndex + state.size))
	}
	// Not on right edge
	if (emptyIndex % state.size != state.size-1) {
		moves = append(moves, state.newSwap(emptyIndex, emptyIndex + 1))
	}
	// Not on left edge
	if (emptyIndex % state.size != 0) {
		moves = append(moves, state.newSwap(emptyIndex, emptyIndex - 1))
	}

	return moves
}

// swaps the two elements in the given indexes for a new state
func (state *State) newSwap(emptyIndex int, newIndex int) *State {
	// we have to create a copy of the board because otherwise they will be linked with a pointer
	boardCopy := make([]int, len(state.board))
	copy(boardCopy, state.board)
	newState := &State{
		board: boardCopy,
		size: state.size,
		complexity: state.complexity+1,
	}
	val := state.board[newIndex]
	newState.board[newIndex] = state.board[emptyIndex]
	newState.board[emptyIndex] = val
	return newState
}

func (state *State) hash() string {
	var bs = make([]byte, len(state.board))
	for i := 0; i < len(state.board); i++ {
		bs[i] = byte(state.board[i])
	}
	return b64.StdEncoding.EncodeToString(bs);
}