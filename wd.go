// implementation of walking distance heurestic
package main

import (
	"math/bits"
)

type visitSt struct {
	cost int8
	board [16]int8
	e int8
}

// creates a new visit struct
func nVisit(cost int8, board [16]int8, e int8) *visitSt {
	return &visitSt{
		cost,
		board,
		e,
	}
}


func code(board [16]int8) int {
	r := 0 
	b := bits.Len(uint(BOARD_ROW_SIZE))
	for i := range board {
		r |= int(board[i]) << (b*i)
	}
	return r
}

type WalkingDistance struct {
	// walking distance table
	table map[int]int8 
	// key is the value of the correct cell
	// "map value int" is the index to the correct value cell
	solved map[int]int
	// BOARD_ROW_SIZE
	size int 
	// BOARD_ROW_SIZE bit length
	bitLength int 
}

func NewWD(rowSize int) *WalkingDistance {
	wd := &WalkingDistance{
		size: rowSize,
		bitLength: bits.Len(uint(rowSize)),
	}
	solvedArr := startingPoint(int16(wd.size))
	ls := make([]int, 16)
	for i := 0; i < 15; i++ {
		ls[i] = int(solvedArr[i])
	}

	// map representation of solved values to used for faster lookup time 
	wd.solved = make(map[int]int)
	for i,val := range ls {
		wd.solved[val] = i 
	}

	wd.GenerateTable()
	return wd 
}

func (wd *WalkingDistance) GenerateTable() *WalkingDistance { 
	wd.table = make(map[int]int8)
	size := int8(wd.size)
	// TODO Fix hard coded solved 
	solved := [16]int8{4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 3}  
	visitable := make(chan *visitSt, 92850)
	visitable <- nVisit(0,solved,size-1)
	count := 0 
	for visit := range visitable {
		key := code(visit.board)
		if _, found := wd.table[key]; found {
			continue 
		}
		wd.table[key] = visit.cost
		for _,d := range []int8{-1,1} {
			if 0 <= (visit.e + d) && (visit.e + d) < size {
				var i int8
				for i < size {
					if visit.board[size*(visit.e+d) + i] > 0 {
						var newBoard [16]int8
						copy(newBoard[:],visit.board[:])
						newBoard[size*(visit.e+d) + i] -= 1
						newBoard[size*visit.e+i] += 1
						visitable <- nVisit(visit.cost+1,newBoard, visit.e+d)
						count++
					}
					i++
				}
			} 
		}
		if count == 92850 {
			close(visitable)
		}
	}
	return wd
}

// calculate the walking distance from the given board
func (wd *WalkingDistance) Calculate(board []t_cell) int {
	heurestic := 0 
	rowCode := 0 
	colCode := 0 

	for i,val := range board {
		if val == 0 {
			continue
		}
		// index that the value should be set to 
		corIdx := wd.solved[int(val)]
		// vertical and horizontal indexs of the current position
		xi, yi := i % wd.size, i / wd.size
		// xCor = vertical index of the correct position
		// yCor = horizontal index of the correct position
		xCor, yCor := corIdx % wd.size, corIdx / wd.size
		// TODO Explain this 
		rowCode += 1 << (wd.bitLength*(wd.size*yi+yCor))
		colCode += 1 << (wd.bitLength*(wd.size*xi+xCor))

		if yCor == yi {
			// calculate vertical heursitic increments
			for yInc := i+1; yInc < i - i%wd.size + wd.size; yInc++ {
				if int(yInc) < len(board) {
					yVal := wd.solved[int(board[int(yInc)])]
					if yVal / wd.size == yi && yVal < corIdx {
						heurestic += 2
					}
				}
			}
		}

		if xCor == xi {
			// calculate horizontal heursitic increments
			for xInc := i + wd.size; xInc < wd.size*wd.size; xInc+=4 {
				if xInc < len(board) {
					kVal := wd.solved[int(board[int(xInc)])]
					if kVal % wd.size == xi && kVal < corIdx {
						heurestic += 2
					}
				}
			}
		}
	}
	heurestic += int(wd.table[rowCode] + wd.table[colCode]) 

	return heurestic
}


// util functions for debugging mostly
func sum(l []int8) int {
	sum := 0
	for _,v := range l {
		sum += int(v) 
	}
	return sum
}

func avg(l []int8) int {
	return sum(l)/len(l)
}

func max(l []int8) int8 {
	var max int8
	for _,v := range l {
		if v > max {
			max = v 
		}
	}
	return max
}