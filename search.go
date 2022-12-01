package main

import "fmt"

// Contains all the important variables for the search
type SearchState struct {
	heuristic       int
	state           *State
	horizontalBoard []int
}

// Transition given vertical board to horizontal representation of the given board
func calculateHorizontalBoard(board []int) []int {
	board = startingPoint(4)
	horizontalBoard := make([]int, len(board))
	copy(horizontalBoard, board)
	// make the board list be a horizontal representation of the puzzle board
	for i := 0; i < int(BOARD_ROW_SIZE); i++ {
		for j := 0; j < int(BOARD_ROW_SIZE); j++ {
			horizontalBoard[i*BOARD_ROW_SIZE+j] = board[i+j*BOARD_ROW_SIZE]
		}
	}

	return horizontalBoard
}

// Evaluate the invert distance of t
func invertDistance(board []int) int {
	// theory of this heuristic evaluation is based on this article https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm

	// Calculate the vertical inversions
	inv := 0
	for i := 0; i < int(BOARD_ROW_SIZE*BOARD_ROW_SIZE); i++ {
		if board[i] != 0 {
			for j := 0; j < i; j++ {
				if board[i] < board[j] {
					inv++
				}
			}
		}
	}
	vertical := inv/3 + 1

	horizontalBoard := calculateHorizontalBoard(board)
	// calculate the horizontal inversions
	inv = 0
	for i := range board {
		// true value of the node so we have to minus one
		value := board[i] - 1
		if value != -1 {
			id := 0
			for j := range horizontalBoard {
				if horizontalBoard[j] == value {
					id = j
					break
				}
			}
			inv += abs(id - i)
		}
	}
	horizontal := inv/3 + 1

	// sum horizontal and vertical distance to get the sum of invert distance
	return vertical + horizontal
}

// node struct of search state, this helps the idastar algorithm to keep track of the depth the algorithm has gone through
type Node struct {
	SearchState
	depth int
}

// create a new search struct from a state
func NewSearch(state *State) *SearchState {
	srh := &SearchState{
		state:     state,
		heuristic: invertDistance(state.board),
	}

	return srh
}

// Iterative Deepening A* search algorithm
func (search *SearchState) IDAStar(maxDepth int) *Node {
	// TODO Figure out what we want to return when the calculations are a success
	// https://en.wikipedia.org/wiki/Iterative_deepening_A*
	depth := 0
	cutoff := search.heuristic
	for depth < maxDepth {
		status, cut, res := IDASearch(search, cutoff, depth)
		if status == SUCCESS {
			return res
		} else if status == CUTOFF {
			cutoff = cut
		}
		depth = res.depth
	}

	return nil
}

// Status constatns of the search state,
type STATUS = int8

// These constants describe what is the current state of the search algorithm
const SUCCESS = STATUS(1)
const CUTOFF = STATUS(2)

// IDASearch
func IDASearch(search *SearchState, cutoff int, startDepth int) (STATUS, int, *Node) {
	h := invertDistance(search.state.board)
	f := h + startDepth
	if f > cutoff {
		return 0, 1, &Node{
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
		// idea if has seen next why bother calculating? we can maybe use cache to store that hash and the state so we don't have to recalculate the whole thing if it has seen the current
		status, probableCutoff, node := IDASearch(&SearchState{
			heuristic: invertDistance(next.board),
			state:     next,
		}, cutoff, startDepth+1)
		if status == CUTOFF {
			stop = true
			nextCutoff = probableCutoff
		} else if status == SUCCESS {
			return status, 0, node
		}
	}
	if stop {
		fmt.Println("returning cutoff")
		return CUTOFF, nextCutoff, nil
	}
	current.heuristic = current.invertDistanceFromMove()
	if current.isSuccess() {
		return SUCCESS, 0, current
	}
	fmt.Println("returning board")
	fmt.Println(current.state.board)
	return -1, 0, current
}

func (search *SearchState) findIndexHorizontal(num int) int {
	if search.horizontalBoard == nil {
		panic("horizontalBoard must be defined if you want to call function findIndexHorizontal")
	}
	for i := range search.horizontalBoard {
		// fmt.Printf("i = %d horizontalBoard[i] = %d  num = %d \n", i, search.horizontalBoard[i], num)
		if search.horizontalBoard[i] == num {
			return i
		}
	}
	return -1
}

func (search *SearchState) isSuccess() bool {
	for i := 0; i < BOARD_ROW_SIZE*BOARD_ROW_SIZE-1; i++ {
		if search.state.board[i] != i+1 {
			return false
		}
	}
	return true
}

// TODO: move this function to state.go
func (search *SearchState) invertDistanceFromMove() int {
	if search.state.move == nil {
		// panic("can't calculate invert distance from a nil move in state. invertMoveDistance should always be called if the state.move value is set")
		return search.heuristic
	}
	count := 0
	/**
		When moving a tile vertically, the total number of inversions can change by only -3, -1, +1, and +3 		- First note that a vertical move will shift the tile 3 positions forward or backwards in our line of tiles.
		- There are two cases to consider, depending on the relative value of the three tiles we've skipped over:
		- Case 1: the three skipped tiles are all smaller (or larger) than the moved tile.
			- Moving the tile will either add or fix three inversions, one for each skipped tile. So, the total number of inversions changes by +3 or -3.
		- Case 2: two of the tiles are larger and other is smaller (or vice versa).
			- In this case, there's going to be a net change of +1 or -1 inversions
	**/

	board := search.state.board
	move := search.state.move
	search.horizontalBoard = calculateHorizontalBoard(board)
	switch move.direction {
	case DIRECTION_UP:
		{
			// Case 1: the three skipped tiles are all smaller (or larger) than the moved tile.
			// - Moving the tile will either add or fix three inversions, one for each skipped tile.
			// So, the total number of inversions changes by +3 or -3.

			// board[move.toIndex] is current empty index
			// board[move.emtpyIndex] is the old value of current empty index, aka the value that was moved
			// move.emtpyIndex is the index of old empty location

			idx := move.toIndex + 1
			tile := board[move.emptyIndex]
			for idx < move.emptyIndex {
				if board[idx] > tile {
					count++
				} else {
					count--
				}
				idx++
			}
		}
	case DIRECTION_DOWN:
		{
			idx := move.toIndex - 1        // 11 index 2
			tile := board[move.emptyIndex] // 6
			for idx > move.emptyIndex {
				if board[idx] > tile { // tätä kutsutaan 3 kertaa indexeillä 2,3,4
					count--
				} else {
					count++
				}
				idx--
			}
		}
	case DIRECTION_RIGHT:
		{
			// this has to be tested
			emptyIndex := search.findIndexHorizontal(0)
			toIndex := search.findIndexHorizontal(board[move.emptyIndex])
			idx := toIndex + 1
			tile := search.horizontalBoard[emptyIndex]
			for idx < emptyIndex {
				if search.horizontalBoard[idx] > tile { // tätä kutsutaan 3 kertaa indexeillä 2,3,4
					count++
				} else {
					count--
				}
				idx++
			}
		}
	case DIRECTION_LEFT:
		{
			emptyIndex := search.findIndexHorizontal(0)
			// fmt.Printf("empty index %d, board move empty index %d \n ", emptyIndex, board[move.emptyIndex])
			// fmt.Println(search.horizontalBoard)
			toIndex := search.findIndexHorizontal(board[move.emptyIndex])
			idx := toIndex - 1
			tile := search.horizontalBoard[emptyIndex]
			for idx > emptyIndex {
				if search.horizontalBoard[idx] > tile { // tätä kutsutaan 3 kertaa indexeillä 2,3,4
					count--
				} else {
					count++
				}
				idx--
			}
		}
	}
	heur := search.heuristic
	if count < 0 {
		heur += abs(count) / 3
	} else {
		heur -= abs(count) / 3
	}

	return heur
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}
