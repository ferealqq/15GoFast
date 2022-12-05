package main

import (
	"fmt"
	"time"
)

// Contains all the important variables for the search
type SearchState struct {
	state               *State
	hasSeen             map[int]*State
	states              []*State
	walkingDistance     func([16]t_cell) int
}

// node struct of search state, this helps the idastar algorithm to keep track of the depth the algorithm has gone through
type Node struct {
	SearchState
	depth t_int
}

// create a new search struct from a state
func NewSearch(state *State) *SearchState {
	wd := NewWD(int(state.size))
	srh := &SearchState{
		state:           state,
		walkingDistance: wd.Calculate,
		hasSeen:         make(map[int]*State),
	}

	return srh
}

type result struct {
	status STATUS
	cutoff t_int
}

// Iterative Deepening A* search algorithm
func (search *SearchState) IDAStar(maxRuntimeMS time.Duration) (*SearchState, STATUS) {
	// TODO Figure out what we want to return when the calculations are a success
	// https://en.wikipedia.org/wiki/Iterative_deepening_A*
	cutoff := t_int(search.walkingDistance(search.state.board))
	search.states = []*State{search.state}

	quitTick := time.NewTicker(maxRuntimeMS * time.Millisecond)
	fmt.Println("solving")
	defer quitTick.Stop()
	for {
		select {
		case <-quitTick.C:
			fmt.Println("time limit exceeded")
			return search, TIME_EXCEEDED
		default:
			status, cut := search.IDASearch(cutoff, t_int(0))
			if status == SUCCESS {
				return search, SUCCESS
			} else if status == CUTOFF {
				cutoff = cut
			} else if status == FAILURE {
				return search, FAILURE
			}
		}
	}
}

// Status constatns of the search state,
type STATUS = int8

// These constants describe what is the current state of the search algorithm
const FAILURE = STATUS(0)
const SUCCESS = STATUS(1)
const CUTOFF = STATUS(2)
const TIME_EXCEEDED = STATUS(3)

// IDASearch, returns STATUS, cutoff, cost
func (search *SearchState) IDASearch(cutoff t_int, cost t_int) (STATUS, t_int) {
	state := search.states[len(search.states)-1]
	h := search.walkingDistance(state.board)
	f := t_int(h) + cost
	if f > cutoff {
		return CUTOFF, f
	}
	var current *State
	stop := false
	nextCutoff := cutoff
	for _, next := range state.GetValidStates() {
		key := code(next.board)
		if _, ok := search.hasSeen[key]; ok {
			continue
		}
		search.states = append(search.states, next)
		search.hasSeen[key] = next
		status, probCut := search.IDASearch(cutoff, cost+1)
		if stop == false || probCut < nextCutoff {
			stop = true
			nextCutoff = probCut
			current = search.states[len(search.states)-1]
		}
		if status == CUTOFF {
			stop = true
			nextCutoff = probCut
		} else if status == SUCCESS {
			return status, 0
		}
		// remove last item from the seen list
		delete(search.hasSeen, key)
		search.states = search.states[:len(search.states)-1]
	}

	if current != nil && current.isSuccess() {
		search.state = current
		// if the states does not contain the latest board add latest board to the states
		if _, ok := search.hasSeen[code(current.board)]; !ok {
			search.states = append(search.states, current)
		}
		return SUCCESS, 0
	}

	if stop {
		return CUTOFF, nextCutoff
	}

	return FAILURE, -1
}

func (search *SearchState) printMoves() {
	for _, s := range search.hasSeen {
		if s.move != nil {
			fmt.Printf(" %s ", s.move.directionString())
		}
	}
	fmt.Printf("\n")
}
