package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func TestPerformance(t *testing.T) {
	flag.Parse()
	if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
					log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
	}else{
		// only run the test when we want to capture the memory usage
		// t.Skip()
	}
	maxRuntimeMS := time.Duration(10600)
	boards := [][16]t_cell{
		{3, 4, 6, 5, 1, 7, 2, 14, 13, 15, 11, 8, 10, 0, 9, 12},
		// {2, 5, 8, 9, 4, 1, 14, 7, 6, 10, 3, 15, 13, 0, 11, 12},
		// {0, 4, 11, 7, 6, 1, 5, 12, 2, 13, 15, 8, 9, 10, 14, 3},
		// {5, 2, 8, 10, 13, 6, 15, 12, 9, 11, 7, 4, 3, 0, 14, 1},
		{3, 6, 8, 7, 5, 0, 9, 2, 1, 4, 14, 15, 13, 10, 12, 11},
	}
	perfList := make(map[int][]time.Duration)
	for _,board := range boards {
		id := code(board)
		for i := 0; i < 3; i++ {
			n := time.Now()
			srh := NewSearch(&State{
				size:       BOARD_ROW_SIZE,
				board: 			board,
				complexity: 0,
			})
			node, _ := srh.IDAStar(maxRuntimeMS)
			dur := time.Since(n)
			fmt.Printf("since n %s \n",dur)
			perfList[id] = append(perfList[id],dur)
			assert.True(t,node.state.isSuccess());
		}
	}
}

func TestPerformanceFast(t *testing.T) {
	flag.Parse()
	if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
					log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			// defer pprof.StopCPUProfile()
	}else{
		// only run the test when we want to capture the memory usage
		// t.Skip()
	}
	maxRuntimeMS := time.Duration(10600)
	boards := [][16]t_cell{
		{3, 4, 6, 5, 1, 7, 2, 14, 13, 15, 11, 8, 10, 0, 9, 12},
		// {2, 5, 8, 9, 4, 1, 14, 7, 6, 10, 3, 15, 13, 0, 11, 12},
		// {0, 4, 11, 7, 6, 1, 5, 12, 2, 13, 15, 8, 9, 10, 14, 3},
		// {5, 2, 8, 10, 13, 6, 15, 12, 9, 11, 7, 4, 3, 0, 14, 1},
		{3, 6, 8, 7, 5, 0, 9, 2, 1, 4, 14, 15, 13, 10, 12, 11},
	}
	perfList := make(map[int][]time.Duration)
	for id,board := range boards {
		for i := 0; i < 60; i++ {
			n := time.Now()
			srh := NewSearch(&State{
				size:       BOARD_ROW_SIZE,
				board: 			board,
				complexity: 0,
			})
			node, _ := srh.IDAStar(maxRuntimeMS)
			dur := time.Since(n)
			fmt.Printf("since n %s \n",dur)
			perfList[id] = append(perfList[id],dur)
			assert.True(t,node.state.isSuccess());
		}
	}
	pprof.StopCPUProfile()
	for i,v := range perfList {
		sum := time.Duration(0)
		for _, j := range v {
			sum += j
		}
		avg := sum/time.Duration(len(v))
		fmt.Printf("board index %d ",i)
		fmt.Printf("average %s \n",avg)
	}
}