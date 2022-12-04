package main

import "fmt"

// Contains all the important variables for the search
type SearchState struct {
	state               *State
	memoInvertDistance  MemoizedFunction[t_cell, []t_cell]
	findIndexHorizontal MemoizedFunction[t_cell, t_cell]
	hasSeen             map[string]*State
	states              []*State
	walkingDistance     func([]t_cell) int
}

// node struct of search state, this helps the idastar algorithm to keep track of the depth the algorithm has gone through
type Node struct {
	SearchState
	depth t_cell
}

// create a new search struct from a state
func NewSearch(state *State) *SearchState {
	wd := NewWD(int(state.size))
	srh := &SearchState{
		state:           state,
		walkingDistance: wd.Calculate,
		hasSeen:         make(map[string]*State),
	}

	return srh
}

// Iterative Deepening A* search algorithm
func (search *SearchState) IDAStar(maxDepth t_cell) *SearchState {
	// TODO Figure out what we want to return when the calculations are a success
	// https://en.wikipedia.org/wiki/Iterative_deepening_A*
	cutoff := t_cell(search.walkingDistance(search.state.board))
	search.states = []*State{search.state}

	// TODO don't use maxDepth let's use max time in milliseconds
	for true {
		status, cut := search.IDASearch(cutoff, t_cell(0))
		if status == SUCCESS {
			return search
		} else if status == CUTOFF {
			cutoff = cut
		} else if status == FAILURE {
			return nil
		}
	}
	fmt.Printf("didn't solve puzzle \n")
	return nil
}

// Status constatns of the search state,
type STATUS = int8

// These constants describe what is the current state of the search algorithm
const FAILURE = STATUS(0)
const SUCCESS = STATUS(1)
const CUTOFF = STATUS(2)

// IDASearch, returns STATUS, cutoff, cost
func (search *SearchState) IDASearch(cutoff t_cell, cost t_cell) (STATUS, t_cell) {
	state := search.states[len(search.states)-1]
	h := search.walkingDistance(state.board)
	f := t_cell(h) + cost
	if f > cutoff {
		return CUTOFF, f
	}
	var current *State
	stop := false
	nextCutoff := cutoff
	for _, next := range state.GetValidStates() {
		key := hash(next.board)
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
