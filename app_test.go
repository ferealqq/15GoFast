package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that the GetValidPositions function returns correct values
func TestGetValidPositions(t *testing.T){
	app := NewApp()
	move := app.GetValidPositions(0)
	assert.Equal(t,Contains(move, 4), true)
	assert.Equal(t,Contains(move, 1), true)

	move = app.GetValidPositions(3)
	assert.Equal(t,Contains(move, 2), true)
	assert.Equal(t,Contains(move, 7), true)

	move = app.GetValidPositions(2)
	assert.Equal(t,Contains(move, 1), true)
	assert.Equal(t,Contains(move, 3), true)
	assert.Equal(t,Contains(move, 6), true)

	move = app.GetValidPositions(5)
	assert.Equal(t,Contains(move, 1), true)
	assert.Equal(t,Contains(move, 4), true)
	assert.Equal(t,Contains(move, 6), true)
	assert.Equal(t,Contains(move, 9), true)

	move = app.GetValidPositions(13)
	assert.Equal(t,Contains(move, 9), true)
	assert.Equal(t,Contains(move, 12), true)
	assert.Equal(t,Contains(move, 14), true)

	move = app.GetValidPositions(12)
	assert.Equal(t,Contains(move, 8), true)
	assert.Equal(t,Contains(move, 13), true)
}

// Test that the ValidateMove function produces the correct results everytime
func TestValidateMove(t *testing.T){
	app := NewApp()

	assert.True(t, app.ValidateMove([2]int{0, 1},0))
	assert.True(t, app.ValidateMove([2]int{0, 4},0))
	assert.False(t, app.ValidateMove([2]int{0, 5},0))

	assert.True(t, app.ValidateMove([2]int{7, 3},7))
	assert.True(t, app.ValidateMove([2]int{7, 6},7))
	assert.True(t, app.ValidateMove([2]int{7, 11},7))

	assert.True(t, app.ValidateMove([2]int{10, 6},10))
	assert.True(t, app.ValidateMove([2]int{10, 9},10))
	assert.True(t, app.ValidateMove([2]int{10, 11},10))
	assert.True(t, app.ValidateMove([2]int{10, 14},10))

	assert.True(t, app.ValidateMove([2]int{10, 6},10))
	assert.True(t, app.ValidateMove([2]int{10, 9},10))
	assert.True(t, app.ValidateMove([2]int{10, 11},10))
	assert.True(t, app.ValidateMove([2]int{10, 14},10))

	assert.True(t, app.ValidateMove([2]int{15, 14},15))
	assert.True(t, app.ValidateMove([2]int{15, 11},15))
	assert.False(t, app.ValidateMove([2]int{15, 10},15))
	assert.False(t, app.ValidateMove([2]int{15, 19},15))
}

// Test that the GetRandomMove function will produce a random yet valid move every time for the player 
func TestGetRandomMove(t *testing.T){
	app := NewApp()
	Times(15, func(index int) []int {
		move := app.GetRandomMove(index)
		assert.True(t,app.ValidateMove(move, index))	
		return nil
	})
	Times(15, func(index int) []int {
		move := app.GetRandomMove(index)
		assert.True(t,app.ValidateMove(move, index))	
		return nil
	})
}