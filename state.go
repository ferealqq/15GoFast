package main

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Describes how many rows the board has
const BOARD_ROW_SIZE = t_cell(4)
type t_cell = int16

type t_direction int8

// Directions
const (
	DIRECTION_UP    = t_direction(0)
	DIRECTION_DOWN  = t_direction(1)
	DIRECTION_LEFT  = t_direction(2)
	DIRECTION_RIGHT = t_direction(3)
)

// Describes what kind of a move has been executed
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
	board      []t_cell
	complexity t_cell
	move       *Move
}

// returns pointer to a new state with clean board.
func NewState() *State {
	return &State{
		size:       BOARD_ROW_SIZE,
		board:      startingPoint(BOARD_ROW_SIZE),
		complexity: 0,
	}
}

// Genereates a State with a board that has shuffeled with N transitions where N is complexity
func GenerateState(complexity t_cell) (*State, error) {
	// https://www.geeksforgeeks.org/check-instance-15-puzzle-solvable/
	state := NewState()
	visited := []string{}
	olds := []*State{}

	rand.Seed(time.Now().UnixNano())
	for state.complexity < complexity {
		visited = append(visited, hash(state.board))
		sts := state.GetValidStates()
		filtered := []*State{}
		for i := range sts {
			for j := range visited {
				if visited[j] != hash(sts[i].board) {
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

func startingPoint(size t_cell) []t_cell {
	res := make([]t_cell, size*size)
	
	for i := t_cell(0); i < (size*size)-1; i++ {
		res[i] = i + 1
	}

	return res
}

// Get the first index of a element in a given array, returns -1 if not found
func GetElementIndex[T comparable](arr []T, element T) int {
	// get empty index
	for i := 0; i < len(arr); i++ {
		if arr[i] == element {
			return i
		}
	}
	return -1
}

// This code has to be very optimized, at the moment this will call too manu mutations

// Calculate all the states that the current state can be mutated to
func (state *State) GetValidStates() []*State {
	// This function is currently O(N^2) - O(N^N), where N = 4
	// appned O(N)
	// newSwap O(N)*
	var states = []*State{}

	// get emtpy index from the board
	emptyIndex := t_cell(GetElementIndex(state.board, 0))

	// the basic logic of puzzle 15 game is that you cannot move off the grid 16x16
	// and you have to move one step at the time at 16x16 grid. Which means that you have to move either +1 -1 +4 -4
	// for example you cannot move up if you are in the top row of the grid.

	// not on the first line
	if emptyIndex-state.size >= 0 {
		states = append(states, state.newSwap(&Move{
			emptyIndex,
			emptyIndex - state.size,
			DIRECTION_UP,
		}))
	}
	// Not on last line
	if emptyIndex+state.size < t_cell(len(state.board)) {
		states = append(states, state.newSwap(&Move{
			emptyIndex,
			emptyIndex + state.size,
			DIRECTION_DOWN,
		}))
	}
	// Not on right edge
	if emptyIndex%state.size != state.size-1 {
		states = append(states, state.newSwap(&Move{
			emptyIndex,
			emptyIndex + 1,
			DIRECTION_RIGHT,
		}))
	}
	// Not on left edge
	if emptyIndex%state.size != 0 {
		states = append(states, state.newSwap(&Move{
			emptyIndex,
			emptyIndex - 1,
			DIRECTION_LEFT,
		}))
	}

	return states
}

// swaps the two elements in the given indexes for a new state
func (state *State) newSwap(move *Move) *State {
	// we have to create a copy of the board because otherwise they will be linked with a pointer
	boardCopy := make([]t_cell, len(state.board))
	copy(boardCopy, state.board)
	newState := &State{
		board:      boardCopy,
		size:       state.size,
		complexity: state.complexity + 1,
		move:       move,
	}
	val := state.board[move.toIndex]
	newState.board[move.toIndex] = state.board[move.emptyIndex]
	newState.board[move.emptyIndex] = val
	return newState
}

// create a unique value for the board represented as string
func hash(board []t_cell) string {
	var bs = make([]byte, len(board))
	for i := 0; i < len(board); i++ {
		bs[i] = byte(board[i])
	}
	return b64.StdEncoding.EncodeToString(bs)
}
