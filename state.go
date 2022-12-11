package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Describes how many rows the board has
const BOARD_ROW_SIZE = t_cell(4)
const BOARD_LENGTH = BOARD_ROW_SIZE * BOARD_ROW_SIZE

type t_int = int16
type t_cell = int8

type t_direction int8

// Directions
const (
	DIRECTION_UP    = t_direction(0)
	DIRECTION_DOWN  = t_direction(1)
	DIRECTION_LEFT  = t_direction(2)
	DIRECTION_RIGHT = t_direction(3)
)

// DEPRICATED Describes what kind of a move has been executed
type Move struct {
	emptyIndex t_cell
	toIndex    t_cell
	direction  t_direction
}

func (m *Move) directionString() string {
	switch m.direction {
	case DIRECTION_DOWN:
		return "Down"
	case DIRECTION_UP:
		return "Up"
	case DIRECTION_LEFT:
		return "Left"
	case DIRECTION_RIGHT:
		return "Right"
	}
	return ""
}

func (m *Move) Print() {
	switch m.direction {
	case DIRECTION_DOWN:
		fmt.Println("Down")
	case DIRECTION_UP:
		fmt.Println("Up")
	case DIRECTION_LEFT:
		fmt.Println("Left")
	case DIRECTION_RIGHT:
		fmt.Println("Right")
	}
}

// State of the 15 puzzle board
type State struct {
	size       t_cell
	board      [16]t_cell
	complexity t_int
}

// returns pointer to a new state with clean board.
func NewState() *State {
	return &State{
		size:       BOARD_ROW_SIZE,
		board:      startingPoint(t_cell(BOARD_ROW_SIZE)),
		complexity: 0,
	}
}

// Genereates a State with a board that has shuffeled with N transitions where N is complexity
func GenerateState(complexity t_int) (*State, error) {
	// https://www.geeksforgeeks.org/check-instance-15-puzzle-solvable/
	state := NewState()
	visited := []int{}
	olds := []*State{}

	rand.Seed(time.Now().UnixNano())
	for state.complexity < complexity {
		visited = append(visited, code(state.board))
		sts := state.GetValidStates()
		filtered := []*State{}
		for _, next := range sts {
			if next == nil {
				continue
			}
			for j := range visited {
				if visited[j] != code(next.board) {
					filtered = append(filtered, next)
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

func startingPoint(size t_cell) [16]t_cell {
	var res [16]t_cell

	for i := t_cell(0); i < (size*size)-1; i++ {
		res[i] = i + 1
	}

	return res
}

// Get the first index of a element in a given array, returns -1 if not found
func getElementIndex(arr [16]t_cell, element t_cell) t_cell {
	for i := 0; i < int(BOARD_LENGTH); i++ {
		if arr[i] == element {
			return t_cell(i)
		}
	}
	return -1
}

// This code has to be very optimized, at the moment this will call too manu mutations

// Calculate all the states that the current state can be mutated to
func (state *State) GetValidStates() [4]*State {
	// This function is currently O(N^2) - O(N^N), where N = 4
	// appned O(N)
	// newSwap O(N)*
	var states [4]*State

	// get emtpy index from the board
	emptyIndex := getElementIndex(state.board, t_cell(0))

	// the basic logic of puzzle 15 game is that you cannot move off the grid 16x16
	// and you have to move one step at the time at 16x16 grid. Which means that you have to move either +1 -1 +4 -4
	// for example you cannot move up if you are in the top row of the grid.

	// not on the first line
	if emptyIndex-state.size >= 0 {
		states[0] = state.newSwap(
			emptyIndex,
			emptyIndex-state.size,
		)

	}
	// Not on last line
	if emptyIndex+state.size < t_cell(len(state.board)) {
		states[1] = state.newSwap(
			emptyIndex,
			emptyIndex+state.size,
		)
	}
	// Not on right edge
	if emptyIndex%state.size != state.size-1 {
		states[2] = state.newSwap(
			emptyIndex,
			emptyIndex+1,
		)

	}
	// Not on left edge
	if emptyIndex%state.size != 0 {
		states[3] = state.newSwap(
			emptyIndex,
			emptyIndex-1,
		)

	}

	return states
}

// swaps the two elements in the given indexes for a new state
func (state *State) newSwap(emptyIndex t_cell, toIndex t_cell) *State {
	// we have to create a copy of the board because otherwise they will be linked with a pointer
	var boardCopy [16]t_cell
	copy(boardCopy[:], state.board[:])
	newState := &State{
		board:      boardCopy,
		size:       state.size,
		complexity: state.complexity + 1,
	}
	val := state.board[toIndex]
	newState.board[toIndex] = state.board[emptyIndex]
	newState.board[emptyIndex] = val
	return newState
}

// depricated should rather use code(board) === code(startingPoint(4))
func (state *State) debugIsSuccess() bool {
	size := t_cell(BOARD_ROW_SIZE*BOARD_ROW_SIZE - 1)
	for i := t_cell(0); i < size; i++ {
		if state.board[i] != i+1 {
			return false
		}
	}
	return true
}
