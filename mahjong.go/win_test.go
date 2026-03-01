package mahjong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_search(t *testing.T) {
	t.Run("hand with multiple winning combinations", func(t *testing.T) {
		tiles := NewTileBag([]Tile{
			TileDots1, TileDots1, TileDots1,
			TileDots2, TileDots2, TileDots2,
			TileDots3, TileDots3, TileDots3,
			TileDragonsWhite, TileDragonsWhite,
		})
		result := search(tiles)
		assert.Equal(t, []Melds{
			{
				{Type: MeldChi, Tiles: []Tile{"13一筒", "14二筒", "15三筒"}},
				{Type: MeldChi, Tiles: []Tile{"13一筒", "14二筒", "15三筒"}},
				{Type: MeldChi, Tiles: []Tile{"13一筒", "14二筒", "15三筒"}},
				{Type: MeldEyes, Tiles: []Tile{"46白板"}},
			},
			{
				{Type: MeldPong, Tiles: []Tile{"13一筒"}},
				{Type: MeldPong, Tiles: []Tile{"14二筒"}},
				{Type: MeldPong, Tiles: []Tile{"15三筒"}},
				{Type: MeldEyes, Tiles: []Tile{"46白板"}},
			},
		}, result)
	})
	t.Run("hand with odd number of tiles", func(t *testing.T) {
		tiles := NewTileBag([]Tile{
			TileDots1, TileDots1, TileDots1,
			TileDots2, TileDots2, TileDots2,
			TileDots3, TileDots3, TileDots3,
			TileDragonsWhite,
		})
		result := search(tiles)
		assert.Empty(t, result)
	})
	t.Run("eyes only", func(t *testing.T) {
		tiles := NewTileBag([]Tile{
			TileDragonsWhite, TileDragonsWhite,
		})
		result := search(tiles)
		assert.Equal(t, []Melds{{{Type: MeldEyes, Tiles: []Tile{"46白板"}}}}, result)
	})
}

func Benchmark_search(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tiles := NewTileBag([]Tile{
			TileDots1, TileDots1, TileDots1,
			TileDots2, TileDots2, TileDots2,
			TileDots3, TileDots3, TileDots3,
			TileDragonsWhite, TileDragonsWhite,
		})
		search(tiles)
	}
}

func Test_score(t *testing.T) {
	t.Run("zi mo ping hu", func(t *testing.T) {
		round := &Round{
			Turn:  0,
			Hands: [4]Hand{{}},
		}
		melds := Melds{
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileCharacters4, TileCharacters5, TileCharacters6}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 4, score(round, 0, melds))
	})
	t.Run("ping hu from discard", func(t *testing.T) {
		round := &Round{
			Turn:  2,
			Hands: [4]Hand{{}},
		}
		melds := Melds{
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileCharacters4, TileCharacters5, TileCharacters6}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 4, score(round, 0, melds))
	})
	t.Run("pong pong hu from discard", func(t *testing.T) {
		round := &Round{
			Turn:  2,
			Hands: [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileCharacters4, TileCharacters4, TileCharacters4}},
			{Type: MeldGang, Tiles: []Tile{TileBamboo2, TileBamboo2, TileBamboo2, TileBamboo2}},
			{Type: MeldPong, Tiles: []Tile{TileBamboo4, TileBamboo4, TileBamboo4}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 2, score(round, 0, melds))
	})
	t.Run("flowers", func(t *testing.T) {
		round := &Round{
			Turn:  2,
			Hands: [4]Hand{{Flowers: []Tile{TileCat, TileGentlemen1, TileGentlemen2}}},
		}
		melds := []Meld{
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileCharacters4, TileCharacters4, TileCharacters4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 2, score(round, 0, melds))
	})
	t.Run("dragons", func(t *testing.T) {
		round := &Round{
			Turn:  2,
			Hands: [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldPong, Tiles: []Tile{TileDragonsRed, TileDragonsRed, TileDragonsRed}},
			{Type: MeldGang, Tiles: []Tile{TileDragonsWhite, TileDragonsWhite, TileDragonsWhite, TileDragonsWhite}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 2, score(round, 0, melds))
	})
	t.Run("seat wind", func(t *testing.T) {
		round := &Round{
			Dealer: 1,
			Turn:   2,
			Hands:  [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileWindsNorth, TileWindsNorth, TileWindsNorth}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 1, score(round, 0, melds))
	})
	t.Run("prevailing wind", func(t *testing.T) {
		round := &Round{
			Dealer: 1,
			Turn:   2,
			Hands:  [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileWindsEast, TileWindsEast, TileWindsEast}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 1, score(round, 0, melds))
	})
	t.Run("seat and prevailing wind", func(t *testing.T) {
		round := &Round{
			Dealer: 0,
			Turn:   2,
			Hands:  [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileWindsEast, TileWindsEast, TileWindsEast}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo4, TileBamboo5, TileBamboo6}},
			{Type: MeldEyes, Tiles: []Tile{TileCharacters9, TileCharacters9}},
		}
		assert.Equal(t, 2, score(round, 0, melds))
	})
	t.Run("full flush", func(t *testing.T) {
		round := &Round{
			Dealer: 0,
			Turn:   2,
			Hands:  [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileDots5, TileDots6, TileDots7}},
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileDots4, TileDots4, TileDots4}},
			{Type: MeldEyes, Tiles: []Tile{TileDots8, TileDots8}},
		}
		assert.Equal(t, 4, score(round, 0, melds))
	})
	t.Run("full flush lesser sequence hand", func(t *testing.T) {
		round := &Round{
			Dealer: 0,
			Turn:   2,
			Hands:  [4]Hand{{Flowers: []Tile{TileGentlemen2}}},
		}
		melds := []Meld{
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileDots5, TileDots6, TileDots7}},
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileDots4, TileDots5, TileDots6}},
			{Type: MeldEyes, Tiles: []Tile{TileDots8, TileDots8}},
		}
		assert.Equal(t, 5, score(round, 0, melds))
	})
	t.Run("half flush", func(t *testing.T) {
		round := &Round{
			Dealer: 0,
			Turn:   2,
			Hands:  [4]Hand{{}},
		}
		melds := []Meld{
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileDots5, TileDots6, TileDots7}},
			{Type: MeldPong, Tiles: []Tile{TileDots1, TileDots1, TileDots1}},
			{Type: MeldPong, Tiles: []Tile{TileDots4, TileDots4, TileDots4}},
			{Type: MeldEyes, Tiles: []Tile{TileWindsWest, TileWindsWest}},
		}
		assert.Equal(t, 2, score(round, 0, melds))
	})
	t.Run("chou ping hu with flowers", func(t *testing.T) {
		round := &Round{
			Dealer: 0,
			Turn:   2,
			Hands: [4]Hand{{
				Flowers: []Tile{TileGentlemen1},
			}},
		}
		melds := []Meld{
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo3, TileBamboo4, TileBamboo5}},
			{Type: MeldChi, Tiles: []Tile{TileBamboo2, TileBamboo3, TileBamboo4}},
			{Type: MeldChi, Tiles: []Tile{TileDots1, TileDots2, TileDots3}},
			{Type: MeldEyes, Tiles: []Tile{TileDots1, TileDots1}},
		}
		assert.Equal(t, 2, score(round, 0, melds))
	})
}

func Test_winnings(t *testing.T) {
	t.Run("default rules", func(t *testing.T) {
		rules := RulesDefault
		actual := winnings(rules, 0, 2, 3)
		expected := [4]int{16, -4, -8, -4}
		assert.Equal(t, expected, actual)
	})
	t.Run("default rules, zi mo", func(t *testing.T) {
		rules := RulesDefault
		actual := winnings(rules, 0, -1, 3)
		expected := [4]int{24, -8, -8, -8}
		assert.Equal(t, expected, actual)
	})
	t.Run("shooter pays", func(t *testing.T) {
		rules := RulesShooter
		actual := winnings(rules, 0, 2, 3)
		expected := [4]int{16, 0, -16, 0}
		assert.Equal(t, expected, actual)
	})
	t.Run("shooter pays, zi mo", func(t *testing.T) {
		rules := RulesShooter
		actual := winnings(rules, 0, -1, 3)
		expected := [4]int{24, -8, -8, -8}
		assert.Equal(t, expected, actual)
	})
	t.Run("default rules, limit", func(t *testing.T) {
		rules := RulesDefault
		actual := winnings(rules, 0, 2, 8)
		expected := [4]int{64, -16, -32, -16}
		assert.Equal(t, expected, actual)
	})
	t.Run("shooter pays, limit", func(t *testing.T) {
		rules := RulesShooter
		actual := winnings(rules, 0, 2, 8)
		expected := [4]int{64, 0, -64, 0}
		assert.Equal(t, expected, actual)
	})
}
