package main

import "fmt"

// Contains all the important variables for the search
type SearchState struct {
	state               *State
	memoInvertDistance  MemoizedFunction[t_cell, []t_cell]
	findIndexHorizontal MemoizedFunction[t_cell, t_cell]
	hasSeen             map[string]*State
	walkingDistance			func([]t_cell) int
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
		state:               state,
		memoInvertDistance:  memoizeBoardCalculation(invertDistance),
		findIndexHorizontal: memoizeBoardCalculation(findIndexInHorizontalBoard),
		// memoize walking distance?
		walkingDistance: 		 wd.Calculate,
		hasSeen: make(map[string]*State),
	}

	return srh
}

// Iterative Deepening A* search algorithm
func (search *SearchState) IDAStar(maxDepth t_cell) *Node {
	// TODO Figure out what we want to return when the calculations are a success
	// https://en.wikipedia.org/wiki/Iterative_deepening_A*
	var depth t_cell = 0 
	// TODO Why does this work but not with walkingDistance?
	cutoff := invertDistance(search.state.board)
	// cutoff := t_cell(search.walkingDistance(search.state.board))

	for depth < maxDepth {
		status, cut, res := search.IDASearch(cutoff, depth)
		if status == SUCCESS {
			return res
		} else if status == CUTOFF {
			cutoff = cut
		}
		depth = res.depth
		// fmt.Printf(" depth < maxDepth, %d < %d \n", depth, maxDepth)
	}

	return nil
}

// Status constatns of the search state,
type STATUS = int8

// These constants describe what is the current state of the search algorithm
const SUCCESS = STATUS(1)
const CUTOFF = STATUS(2)

// IDASearch
func (search *SearchState) IDASearch(cutoff t_cell, startDepth t_cell) (STATUS, t_cell, *Node) {
	h := search.walkingDistance(search.state.board)
	f := t_cell(h) + startDepth
	// fmt.Printf("f %d cutoff %d \n",f,cutoff)
	if f > cutoff {
		return 0, f, &Node{
			*search,
			startDepth,
		}
	}
	current := &Node{
		*search,
		startDepth,
	}
	stop := false
	nextCutoff := cutoff
	sts := search.state.GetValidStates()
	for i := range sts {
		next := sts[i]
		key := hash(next.board)
		if _, ok := search.hasSeen[key]; ok  {
			continue
		}
		// idea if has seen next why bother calculating? we can maybe use cache to store that hash and the state so we don't have to recalculate the whole thing if it has seen the current
		search.state = next
		search.hasSeen[key] = next
		status, _, node := search.IDASearch(cutoff, startDepth+1)
		// if nextH := t_cell(node.walkingDistance(node.state.board)); nextH < nextCutoff { 
		// 	stop = true 
		// 	nextCutoff = nextH
		// 	current = node 
		// }
		if status == CUTOFF {
			stop = true
			// kun täs kutsuu nextiä niin jos search.state muuttuu niin mutatoiko se myös nextiä?
			nextCutoff = t_cell(search.walkingDistance(next.board))
			// fmt.Printf("node depth %d \n",node.depth)
			// current = node
		} else if status == SUCCESS {
			// fmt.Println("rekursiivinen palauttaa")
			return status, 0, node
		}
		// remove last item from the seen list
		delete(search.hasSeen, key);
	}
	if stop {
		// fmt.Println("stop palauttaa")
		return CUTOFF, nextCutoff, current
	}
	if current.isSuccess() {
		// fmt.Println("isSuccess palauttaa")
		return SUCCESS, 0, current
	}
	// fmt.Println("normi palautusta vaa hinkataan")
	return -1, nextCutoff, current
}

func (search *SearchState) printMoves() {
	for _, s := range search.hasSeen {
		if s.move != nil {
			fmt.Printf(" %s ", s.move.directionString())
		}
	}
	fmt.Printf("\n")
}

func (search *SearchState) isSuccess() bool {
	for i := t_cell(0); i < BOARD_ROW_SIZE*BOARD_ROW_SIZE-1; i++ {
		if search.state.board[i] != i+1 {
			return false
		}
	}
	return true
}
