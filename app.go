package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)


const DEFAULT_COMPLEXITY = 100
const DEFAULT_MAX_RUNTIME = 1300
// App struct
type App struct {
	ctx context.Context
	search *SearchState
	complexity t_cell
	maxRuntime time.Duration
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
	app.complexity = DEFAULT_COMPLEXITY
	app.maxRuntime = time.Duration(DEFAULT_MAX_RUNTIME)
	st,err := GenerateState(app.complexity)
	if err != nil {
		// TODO  ?
		panic(err)
	}
	app.search = NewSearch(st)
}
// get current board in the app
func (app *App) GetBoard() []t_cell {
	return app.search.state.board
}
// generate new board with the app complexity 
func (app *App) GenerateBoard() []t_cell {
	st,err := GenerateState(app.complexity)
	if err != nil {
		// TODO  ?
		panic(err)
	}
	app.search = NewSearch(st)
	return app.GenerateBoard()
}

func (app *App) SetComplexity(comp t_cell) {
	app.complexity = comp
}

func (app *App) SetMaxRuntime(milliseconds int) {
	app.maxRuntime = time.Duration(milliseconds)
}

type SolveResult struct {
	Status STATUS
	Iterations [][]t_cell
}

// returns every iteration of the board while solving  
func (app *App) Solve() SolveResult {
	res,status := app.search.IDAStar(app.maxRuntime)
	if status == SUCCESS {
		boards := make([][]t_cell, len(res.states))
		for i,state := range res.states {
			boards[i] = state.board
		}
		return SolveResult{
			SUCCESS,
			boards,
		}
	}else{
		return SolveResult{Status: status}
	}
}

// Greet returns a greeting for the given name
func (app *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Indexs of the borders of the puzzle games grid
var verticalTop = []int{0, 1, 2, 3}
var horizontalLeft = []int{0, 4, 8, 12}
var horizontalRight = []int{3, 7, 11, 15}
var verticalBottom = []int{12, 13, 14, 15}

// Get positions that the player can move to
func (app *App) GetValidPositions(emptyIndex int) []int {
	var moves = []int{}

	// the basic logic of puzzle 15 game is that you cannot move off the grid 16x16
	// and you have to move one step at the time at 16x16 grid. Which means that you have to move either +1 -1 +4 -4
	// for example you cannot move up if you are in the top row of the grid.
	if !Contains(verticalBottom, emptyIndex) {
		moves = append(moves, emptyIndex+4)
	}
	if !Contains(verticalTop, emptyIndex) {
		moves = append(moves, emptyIndex-4)
	}
	if !Contains(horizontalLeft, emptyIndex) {
		moves = append(moves, emptyIndex-1)
	}
	if !Contains(horizontalRight, emptyIndex) {
		moves = append(moves, emptyIndex+1)
	}

	return moves
}

// Checks wheter a certain variable exist in the given array
func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

// Do a simple process count times with iteratee funcion
func Times[T any](count int, iteratee func(index int) T) []T {
	result := make([]T, count)

	for i := 0; i < count; i++ {
		result[i] = iteratee(i)
	}

	return result
}

// checks if the given move is correct or incorrect
func (app *App) ValidateMove(move [2]int, emptyIndex int) bool {
	pos := app.GetValidPositions(emptyIndex)
	var result = false

	for _, v := range pos {
		if move == [2]int{emptyIndex, v} {
			result = true
		}
	}

	return result
}

// return a random move that the player is able to produce
func (app *App) GetRandomMove(emptyIndex int) [2]int {
	moves := app.GetValidPositions(emptyIndex)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })
	return [2]int{emptyIndex, moves[0]}
}