package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInversionDistance(t *testing.T){
	board := []int{1,2,4,8,9,5,10,3,7,14,6,12,13,0,11,15}
	state := NewSearch(NewState())
	inv1 := state.invertDistance(board)
	// https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru/stp_wd_translation_en.htm
	// TODO Not sure that the math behind this is correct, might have to check later
	assert.Equal(t, inv1, 27)
}