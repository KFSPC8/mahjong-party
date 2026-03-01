package mahjong

// MeldType represents the type of a melded set.
type MeldType int

// Allowed meld types.
const (
	MeldChi MeldType = iota
	MeldPong
	MeldGang
	MeldEyes
)

// Meld represents a melded set.
type Meld struct {
	Type  MeldType `json:"type"`
	Tiles []Tile   `json:"tiles"`
}

type Melds []Meld

func (m Melds) Len() int {
	return len(m)
}

func (m Melds) Less(i, j int) bool {
	if m[i].Type < m[j].Type {
		return true
	}
	return m[i].Tiles[0] < m[j].Tiles[0]
}

func (m Melds) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Melds) Tiles() []Tile {
	var tiles []Tile
	for _, meld := range m {
		switch meld.Type {
		case MeldChi:
			tiles = append(tiles, meld.Tiles...)
		case MeldPong:
			tiles = append(tiles, meld.Tiles[0], meld.Tiles[0], meld.Tiles[0])
		case MeldGang:
			tiles = append(tiles, meld.Tiles[0], meld.Tiles[0], meld.Tiles[0], meld.Tiles[0])
		case MeldEyes:
			tiles = append(tiles, meld.Tiles[0], meld.Tiles[0])
		}
	}
	return tiles
}

type TileBag map[Tile]int

func (b TileBag) Cardinality() int {
	c := 0
	for _, count := range b {
		c += count
	}
	return c
}

func (b TileBag) Contains(tile Tile) bool {
	return b[tile] > 0
}

func (b TileBag) Count(tile Tile) int {
	return b[tile]
}

func (b TileBag) Add(tiles ...Tile) {
	for _, tile := range tiles {
		b[tile]++
	}
}

func (b TileBag) Remove(tiles ...Tile) {
	for _, tile := range tiles {
		if b[tile] == 0 {
			continue
		}
		if b[tile] == 1 {
			delete(b, tile)
		} else {
			b[tile]--
		}
	}
}

func (b TileBag) RemoveN(tile Tile, n int) bool {
	if b[tile] < n {
		return false
	}
	if b[tile] == n {
		delete(b, tile)
	} else {
		b[tile] -= n
	}
	return true
}

func NewTileBag(tiles []Tile) TileBag {
	b := TileBag{}
	for _, tile := range tiles {
		b[tile]++
	}
	return b
}

// Hand represents all the tiles belonging to a player.
type Hand struct {
	Flowers   []Tile  `json:"flowers"`
	Revealed  Melds   `json:"revealed"`
	Concealed TileBag `json:"concealed,omitempty"`
	Finished  []Tile  `json:"finished,omitempty"`
}

// View returns another player's view of a hand.
func (h Hand) View() Hand {
	return Hand{
		Flowers:   h.Flowers,
		Revealed:  h.Revealed,
		Concealed: TileBag{"": h.Concealed.Cardinality()},
		Finished:  h.Finished,
	}
}

// Direction represents a wind direction.
type Direction int

// All possible directions.
const (
	DirectionEast Direction = iota
	DirectionSouth
	DirectionWest
	DirectionNorth
)

// Phase constrains what actions are currently possible.
type Phase string

const (
	// PhaseDraw represents the draw phase, when the player whose turn it
	// currently is may draw a tile or chi the last discarded tile, and any
	// player may pong the last discarded tile.
	PhaseDraw Phase = "draw"

	// PhaseDiscard represents the discard phase, when the player whose turn it
	// is has 14 tiles in their hand and must discard a tile. They may also
	// reveal a concealed gang or win by self-draw.
	PhaseDiscard Phase = "discard"
)

// Result represents the outcome of a round.
type Result struct {
	// Dealer is the integer offset of the dealer for the round.
	Dealer int `json:"dealer"`

	// Wind is the prevailing wind for the round.
	Wind Direction `json:"wind"`

	// Winner is the integer offset of the winner for the round, or -1 if the round ended in a draw.
	Winner int `json:"winner"`

	// Loser is the integer offset of the player who threw the winning tile, or -1 if the winner won by zi mo.
	Loser int `json:"loser"`

	// Points is how much the winning hand was worth.
	Points int `json:"points"`

	// WinningTiles is the set of flowers and tiles belonging to the winner.
	WinningTiles []Tile `json:"winning_tiles"`
}
