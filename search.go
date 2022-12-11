package main

import (
	"fmt"
	"time"
)

// Contains all the important variables for the search
type SearchState struct {
	state       *State
	hasSeen     map[int]*State
	Heuristic   heurFn
	successCode int
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
		state:       state,
		hasSeen:     make(map[int]*State),
		successCode: codeUniq(startingPoint(state.size)),
	}

	srh.Heuristic = srh.memo(wd.Calculate)
	// srh.Heuristic = wd.Calculate

	return srh
}

type heurFn = func([16]t_cell) int

func (search *SearchState) memo(fn heurFn) heurFn {
	cache := make(map[int]int)

	return func(input [16]t_cell) int {
		key := codeUniq(input)
		if val, found := cache[key]; found {
			return val
		}
		val := fn(input)
		cache[key] = fn(input)
		return val
	}
}

type result struct {
	status STATUS
	cutoff t_int
}

// Iterative Deepening A* search algorithm
func (search *SearchState) IDAStar(maxRuntimeMS time.Duration) (*SearchState, STATUS) {
	// https://en.wikipedia.org/wiki/Iterative_deepening_A*
	cutoff := t_int(search.Heuristic(search.state.board))
	quitTick := time.NewTicker(maxRuntimeMS * time.Millisecond)
	// miten tähän pitäis reagoida? chekkaatko onko se tässä tilassa ja heitä feilu vai mitä? ja mites muuten generate state voi mennä tohon tilaan? Ton ei pitäis olla matemaattisesti mahdollista mennä tohon tilaan.
	// [16]int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 10, 12, 13, 15, 14, 0}
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
	state := search.state
	h := search.Heuristic(state.board)
	f := t_int(h) + cost
	if f > cutoff {
		return CUTOFF, f
	}
	var old State = *search.state
	stop := false
	nextCutoff := cutoff
	for _, next := range state.GetValidStates() {
		if next == nil {
			continue
		}
		key := codeUniq(next.board)
		// checks if the board is in the starting position. successCode === startingPosition codeUniq
		if _, ok := search.hasSeen[key]; ok {
			continue
		}
		search.hasSeen[key] = next
		search.state = next
		if key == search.successCode {
			return SUCCESS, 0
		}

		status, probCut := search.IDASearch(cutoff, cost+1)
		if stop == false || probCut < nextCutoff {
			stop = true
			nextCutoff = probCut
		}
		if status == CUTOFF {
			stop = true
			nextCutoff = probCut
		} else if status == SUCCESS {
			return status, 0
		}
		// remove last item from the seen list
		delete(search.hasSeen, key)
		search.state = &old
	}

	if stop {
		return CUTOFF, nextCutoff
	}

	return FAILURE, -1
}