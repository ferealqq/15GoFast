package main

import (
	"context"
	"fmt"

	_ "net/http/pprof"
	"time"
)

const DEFAULT_COMPLEXITY = 130
const DEFAULT_MAX_RUNTIME = 1300

// App struct
type App struct {
	ctx        context.Context
	search     *SearchState
	complexity t_int
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
	st, err := GenerateState(app.complexity)
	if err != nil {
		// TODO  ?
		panic(err)
	}
	app.search = NewSearch(st)
}

// get current board in the app
func (app *App) GetBoard() [16]t_cell {
	return app.search.state.board
}

// generate new board with the app complexity
func (app *App) GenerateBoard() [16]t_cell {
	st, err := GenerateState(app.complexity)
	if err != nil {
		// TODO  ?
		panic(err)
	}
	app.search = NewSearch(st)
	return app.GetBoard()
}

func (app *App) SetComplexity(comp t_int) {
	app.complexity = comp
}

func (app *App) SetMaxRuntime(milliseconds int) {
	app.maxRuntime = time.Duration(milliseconds)
}

type IterationData struct {
	Move struct {
		EmptyIndex t_cell
		ToIndex    t_cell
		Direction  t_direction
	}
	Board [16]t_cell
}

type SolveResult struct {
	Status         STATUS
	Iterations     []IterationData
	IterationCount int
	//time elapsed in milliseconds
	TimeElapsed time.Duration
}

// returns every iteration of the board while solving
func (app *App) Solve() SolveResult {
	now := time.Now()
	res, status := app.search.IDAStar(app.maxRuntime)
	elapsed := time.Since(now)
	if status == SUCCESS {
		boards := make([]IterationData, len(res.states))
		for i, state := range res.states {
			boards[i] = IterationData{Board: state.board, Move: struct {
				EmptyIndex t_cell
				ToIndex    t_cell
				Direction  t_direction
			}{
				EmptyIndex: state.move.emptyIndex,
				ToIndex: state.move.toIndex,
				Direction: state.move.direction,
			}}
			fmt.Println(boards[i])
		}
		ms := elapsed / time.Millisecond
		fmt.Println(len(boards))
		return SolveResult{
			SUCCESS,
			boards,
			len(boards),
			ms,
		}
	} else {
		return SolveResult{Status: status}
	}
}