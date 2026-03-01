package mahjong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isValidSequence(t *testing.T) {
	tile0, tile1, tile2 := TileBamboo4, TileBamboo2, TileBamboo3
	assert.True(t, isValidSequence(tile0, tile1, tile2))
}
