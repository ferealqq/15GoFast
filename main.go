// Documentation of the 15Puzzle solver.
// This documentation contains all public varibles used in the application
package main

import (
	"embed"
	"log"
	"net/http"

	// _ "net/http/pprof"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "15GoFast",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}


// func main() {
// 	maxRuntimeMS := time.Duration(1500)
// 	st, _ := GenerateState(70)
// 	board := st.board
// 	state := NewState()
// 	state.board = board
// 	fmt.Println("Started sovling board")
// 	fmt.Println(board)
// 	srh := NewSearch(state)
// 	node, _ := srh.IDAStar(maxRuntimeMS)
// 	if !node.state.isSuccess() {
// 		fmt.Println("Failed solving board")
// 	}else{
// 		fmt.Println("Solved")
// 	}

// 	maxRuntimeMS = time.Duration(2600)
// 	st, _ = GenerateState(135)
// 	board = st.board
// 	state = NewState()
// 	state.board = board
// 	fmt.Println("Started sovling board")
// 	fmt.Println(board)
// 	srh = NewSearch(state)
// 	node, _ = srh.IDAStar(maxRuntimeMS)
// 	if !node.state.isSuccess() {
// 		fmt.Println("Failed solving board")
// 	}else{
// 		fmt.Println("Solved")
// 	}

// 	log.Println(http.ListenAndServe("localhost:6060", nil))
// }