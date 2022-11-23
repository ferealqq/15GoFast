package main

// Evaluate the invert distance of t
func invertDistance(board []int) int {
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
	// calculate the horizontal invert distance
	inv = 0
	for i := range board {
		// true value of the node so we have to minus one
		value := board[i] - 1;
		if value != -1 {
			id := 0;
			for j := range horizontalBoard {
				if horizontalBoard[j] == value {
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

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}


