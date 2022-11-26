package main

type SearchState struct {
	heuristic int
	state *State
	horizontalBoard []int
} 

// type Search struct {
// }

func calculateHorizontalBoard(board []int) []int {
	// this horizontal move calculation can be saved in the app state so it wouldn't be called all the time
	horizontalBoard := make([]int, len(board))
	copy(horizontalBoard,board)
	// make the board list be a horizontal representation of the puzzle board
	idx := 0
	for i := 0; i < int(BOARD_ROW_SIZE); i++ {
		for j := 0; j < int(BOARD_ROW_SIZE); j++ {
			horizontalBoard[idx] = j*int(BOARD_ROW_SIZE) + i
			idx++
		}
	}

	return horizontalBoard
}

// Evaluate the invert distance of t
func (search *SearchState) invertDistance(board []int) int {
	// theory of this heuristic evaluation is based on this article https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm

	// Calculate the vertical invert distance 
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

	// // this horizontal move calculation can be saved in the app state so it wouldn't be called all the time
	// horizontalBoard := make([]int, len(board))
	// copy(horizontalBoard,board)
	// // make the board list be a horizontal representation of the puzzle board
	// idx := 0
	// for i := 0; i < int(BOARD_ROW_SIZE); i++ {
	// 	for j := 0; j < int(BOARD_ROW_SIZE); j++ {
	// 		horizontalBoard[idx] = j*int(BOARD_ROW_SIZE) + i
	// 		idx++
	// 	}
	// }

	// calculate the horizontal invert distance
	inv = 0
	for i := range board {
		// true value of the node so we have to minus one
		value := board[i] - 1;
		if value != -1 {
			id := 0;
			for j := range search.horizontalBoard {
				if search.horizontalBoard[j] == value {
					id = j
					break
				}
			}	
			inv += abs(id - i);
		}
	}
	horizontal := inv/3 + 1

	// sum horizontal and vertical distance to get the sum of invert distance
	return vertical + horizontal
}

type Node struct {
	SearchState
	depth int
}

func NewSearch(state *State) *SearchState {
	srh := &SearchState{
		state: state,
		horizontalBoard: calculateHorizontalBoard(state.board),
	}

	srh.heuristic = srh.invertDistance(state.board)

	return srh
}
// Iterative Deepening A* search algorithm
func (search *SearchState) IDAStar(maxDepth int) *Node {
	// TODO Figure out what we want to return when the calculations are a success
	// https://en.wikipedia.org/wiki/Iterative_deepening_A*
	depth := 0
	cutoff := search.heuristic
	for depth < maxDepth {
		suc,cut,res := IDASearch(search, cutoff, depth)
		if suc == SUCCESS {
			return res
		}
		depth = res.depth
		cutoff = cut
	}

	return nil
}

type success = int8;

const SUCCESS = int8(1);

func IDASearch(search *SearchState, cutoff int, startDepth int) (success,int, *Node) {
	h := search.invertDistance(search.state.board)
	f := h + startDepth
	if f > cutoff {
		return 0,1,&Node{
			*search,
			startDepth,
		}
	}
	current := &Node{
		*search,
		startDepth,
	}
	sts := search.state.GetValidStates()
	for i := range sts {
		next := sts[i]
		// idea if has seen next why bother calculating? 
		if nextH := search.invertDistance(next.board); current.heuristic > nextH {
			current = &Node{
				SearchState{
					heuristic: nextH,
					state: next,
				},
				current.depth+1,
			}
			success,cutoff,res := IDASearch(&current.SearchState,current.heuristic,current.depth+1)
			if nextH = res.heuristic; current.heuristic > nextH {
				current = res
			}
			// TODO
			if cutoff == 1 {
				// TODO 
			}
			if success == SUCCESS {
				return 1,0,res
			}
		}
	}
	current.heuristic = current.state.invertDistanceFromMove()
	return 0,0,current
}

// TODO: move this function to state.go
func (state *State) invertDistanceFromMove() int {
	if state.move == nil {
		panic("can't calculate invert distance from a nil move in state. invertMoveDistance should always be called if the state.move value is set")
	}
	/**
		When moving a tile vertically, the total number of inversions can change by only -3, -1, +1, and +3 		- First note that a vertical move will shift the tile 3 positions forward or backwards in our line of tiles. 
		- There are two cases to consider, depending on the relative value of the three tiles we've skipped over: 
		- Case 1: the three skipped tiles are all smaller (or larger) than the moved tile. 
			- Moving the tile will either add or fix three inversions, one for each skipped tile. So, the total number of inversions changes by +3 or -3. 
		- Case 2: two of the tiles are larger and other is smaller (or vice versa). 
			- In this case, there's going to be a net change of +1 or -1 inversions
	**/
	if state.move.direction == DIRECTION_LEFT {

	}

	return 0;
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}


