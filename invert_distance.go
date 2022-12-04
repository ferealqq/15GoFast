package main

// TODO: move this function to state.go
func (search *SearchState) invertDistanceFromMove() t_cell {
	// TODO fix this function
	if search.state.move == nil {
		// panic("can't calculate invert distance from a nil move in state. invertMoveDistance should always be called if the state.move value is set")
		return 0
	}
	var count int = 0
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
			tile := memoHorizontalBoard[emptyIndex]
			for idx < emptyIndex {
				if memoHorizontalBoard[idx] > tile { // tätä kutsutaan 3 kertaa indexeillä 2,3,4
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
			toIndex := search.findIndexHorizontal(board[move.emptyIndex])
			idx := toIndex - 1
			tile := memoHorizontalBoard[emptyIndex]
			for idx > emptyIndex {
				if memoHorizontalBoard[idx] > tile { // tätä kutsutaan 3 kertaa indexeillä 2,3,4
					count--
				} else {
					count++
				}
				idx--
			}
		}
	}
	heur := t_cell(0)
	if count < 0 {
		heur += t_cell(abs(count) / 3)
	} else {
		heur -= t_cell(abs(count) / 3)
	}

	return heur
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// Evaluate the invert distance of t
func invertDistance(board []t_cell) t_cell {
	// theory of this heuristic evaluation is based on this article https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm

	// Calculate the vertical inversions
	inv := 0
	for i := t_cell(0); i < BOARD_ROW_SIZE*BOARD_ROW_SIZE; i++ {
		if board[i] != 0 {
			for j := t_cell(0); j < i; j++ {
				if board[i] < board[j] {
					inv++
				}
			}
		}
	}
	vertical := t_cell(inv/3 + 1)

	// calculate the horizontal inversions
	inv = 0
	for i := range board {
		// true value of the node so we have to minus one
		value := board[i]
		if value != -1 {
			id := 0
			for j := range memoHorizontalBoard {
				if memoHorizontalBoard[j] == t_cell(value) {
					id = j
					break
				}
			}
			inv += abs(id - i)
		}
	}
	horizontal := t_cell(inv/3 + 1)

	// sum horizontal and vertical distance to get the sum of invert distance
	return vertical + horizontal
}

type MemoizedFunction[T interface{}, R interface{}] func(R) T

// because we are doing alot of expensive or semi expensive calculations. It's prefered that we memoize the values that these calls return rather than we call the operations with same values again and again
func memoizeBoardCalculation[T interface{}, R interface{}](fn MemoizedFunction[T, R]) MemoizedFunction[T, R] {
	cache := make(map[interface{}]T)

	return func(input R) T {
		if val, found := cache[input]; found {
			return val
		}
		val := fn(input)
		cache[input] = fn(input)
		return val
	}
}

// Transition given vertical board to horizontal representation of the given board
func calculateHorizontalBoard(rowSize t_cell) []t_cell {
	board := startingPoint(rowSize)
	horizontalBoard := make([]t_cell, len(board))
	copy(horizontalBoard, board)
	// make the board list be a horizontal representation of the puzzle board
	for i := t_cell(0); i < BOARD_ROW_SIZE; i++ {
		for j := t_cell(0); j < BOARD_ROW_SIZE; j++ {
			horizontalBoard[i*BOARD_ROW_SIZE+j] = board[i+j*BOARD_ROW_SIZE]
		}
	}

	return horizontalBoard
}

var memoHorizontalBoard = calculateHorizontalBoard(BOARD_ROW_SIZE)

func findIndexInHorizontalBoard(num t_cell) t_cell {
	for i := range memoHorizontalBoard {
		if memoHorizontalBoard[i] == num {
			return t_cell(i)
		}
	}
	return -1
}
