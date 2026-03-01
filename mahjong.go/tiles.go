package mahjong

// Tile represents a mahjong tile.
type Tile string

// All tiles.
const (
	TileCat          Tile = "01猫"
	TileRat          Tile = "02老鼠"
	TileRooster      Tile = "03公鸡"
	TileCentipede    Tile = "04蜈蚣"
	TileGentlemen1   Tile = "05梅"
	TileGentlemen2   Tile = "06兰"
	TileGentlemen3   Tile = "07菊"
	TileGentlemen4   Tile = "08竹"
	TileSeasons1     Tile = "09春"
	TileSeasons2     Tile = "10夏"
	TileSeasons3     Tile = "11秋"
	TileSeasons4     Tile = "12冬"
	TileDots1        Tile = "13一筒"
	TileDots2        Tile = "14二筒"
	TileDots3        Tile = "15三筒"
	TileDots4        Tile = "16四筒"
	TileDots5        Tile = "17五筒"
	TileDots6        Tile = "18六筒"
	TileDots7        Tile = "19七筒"
	TileDots8        Tile = "20八筒"
	TileDots9        Tile = "21九筒"
	TileBamboo1      Tile = "22一索"
	TileBamboo2      Tile = "23二索"
	TileBamboo3      Tile = "24三索"
	TileBamboo4      Tile = "25四索"
	TileBamboo5      Tile = "26五索"
	TileBamboo6      Tile = "27六索"
	TileBamboo7      Tile = "28七索"
	TileBamboo8      Tile = "29八索"
	TileBamboo9      Tile = "30九索"
	TileCharacters1  Tile = "31一万"
	TileCharacters2  Tile = "32二万"
	TileCharacters3  Tile = "33三万"
	TileCharacters4  Tile = "34四万"
	TileCharacters5  Tile = "35五万"
	TileCharacters6  Tile = "36六万"
	TileCharacters7  Tile = "37七万"
	TileCharacters8  Tile = "38八万"
	TileCharacters9  Tile = "39九万"
	TileWindsEast    Tile = "40东风"
	TileWindsSouth   Tile = "41南风"
	TileWindsWest    Tile = "42西风"
	TileWindsNorth   Tile = "43北风"
	TileDragonsRed   Tile = "44红中"
	TileDragonsGreen Tile = "45青发"
	TileDragonsWhite Tile = "46白板"
)

type Suit int

const (
	SuitInvalid Suit = iota
	SuitFlowers
	SuitDots
	SuitBamboo
	SuitCharacters
	SuitDragons
	SuitWinds
)

var (
	flowerTiles = []Tile{TileCat, TileRat, TileRooster, TileCentipede, TileGentlemen1, TileGentlemen2, TileGentlemen3, TileGentlemen4, TileSeasons1, TileSeasons2, TileSeasons3, TileSeasons4}
	suitedTiles = []Tile{
		TileDots1, TileDots2, TileDots3, TileDots4, TileDots5, TileDots6, TileDots7, TileDots8, TileDots9,
		TileBamboo1, TileBamboo2, TileBamboo3, TileBamboo4, TileBamboo5, TileBamboo6, TileBamboo7, TileBamboo8, TileBamboo9,
		TileCharacters1, TileCharacters2, TileCharacters3, TileCharacters4, TileCharacters5, TileCharacters6, TileCharacters7, TileCharacters8, TileCharacters9,
		TileWindsEast, TileWindsSouth, TileWindsWest, TileWindsNorth,
		TileDragonsRed, TileDragonsGreen, TileDragonsWhite,
	}
)

type FlowerGroup struct {
	Flowers []Tile
	Payout  int
}

var (
	bites = map[Tile][]FlowerGroup{
		TileCat: {
			{Flowers: []Tile{TileCat, TileRat}, Payout: 2},
			{Flowers: []Tile{TileCat, TileRat, TileRooster, TileCentipede}, Payout: 4},
		},
		TileRat: {
			{Flowers: []Tile{TileCat, TileRat}, Payout: 2},
			{Flowers: []Tile{TileCat, TileRat, TileRooster, TileCentipede}, Payout: 4},
		},
		TileRooster: {
			{Flowers: []Tile{TileRooster, TileCentipede}, Payout: 2},
			{Flowers: []Tile{TileCat, TileRat, TileRooster, TileCentipede}, Payout: 4},
		},
		TileCentipede: {
			{Flowers: []Tile{TileRooster, TileCentipede}, Payout: 2},
			{Flowers: []Tile{TileCat, TileRat, TileRooster, TileCentipede}, Payout: 4},
		},
	}
)

func (t Tile) Suit() Suit {
	switch {
	case t == TileCat || t == TileRat || t == TileRooster || t == TileCentipede || t == TileGentlemen1 || t == TileGentlemen2 || t == TileGentlemen3 || t == TileGentlemen4 || t == TileSeasons1 || t == TileSeasons2 || t == TileSeasons3 || t == TileSeasons4:
		return SuitFlowers
	case t == TileDots1 || t == TileDots2 || t == TileDots3 || t == TileDots4 || t == TileDots5 || t == TileDots6 || t == TileDots7 || t == TileDots8 || t == TileDots9:
		return SuitDots
	case t == TileBamboo1 || t == TileBamboo2 || t == TileBamboo3 || t == TileBamboo4 || t == TileBamboo5 || t == TileBamboo6 || t == TileBamboo7 || t == TileBamboo8 || t == TileBamboo9:
		return SuitBamboo
	case t == TileCharacters1 || t == TileCharacters2 || t == TileCharacters3 || t == TileCharacters4 || t == TileCharacters5 || t == TileCharacters6 || t == TileCharacters7 || t == TileCharacters8 || t == TileCharacters9:
		return SuitCharacters
	case t == TileWindsEast || t == TileWindsSouth || t == TileWindsWest || t == TileWindsNorth:
		return SuitWinds
	case t == TileDragonsRed || t == TileDragonsGreen || t == TileDragonsWhite:
		return SuitDragons
	}
	return 0
}

func isFlower(tile Tile) bool {
	return tile.Suit() == SuitFlowers
}

// sequences is a map of tiles to valid tiles for completing a sequence.
var sequences = map[Tile][][2]Tile{
	TileDots1:       {{TileDots2, TileDots3}},
	TileDots2:       {{TileDots1, TileDots3}, {TileDots3, TileDots4}},
	TileDots3:       {{TileDots1, TileDots2}, {TileDots2, TileDots4}, {TileDots4, TileDots5}},
	TileDots4:       {{TileDots2, TileDots3}, {TileDots3, TileDots5}, {TileDots5, TileDots6}},
	TileDots5:       {{TileDots3, TileDots4}, {TileDots4, TileDots6}, {TileDots6, TileDots7}},
	TileDots6:       {{TileDots4, TileDots5}, {TileDots5, TileDots7}, {TileDots7, TileDots8}},
	TileDots7:       {{TileDots5, TileDots6}, {TileDots6, TileDots8}, {TileDots8, TileDots9}},
	TileDots8:       {{TileDots6, TileDots7}, {TileDots7, TileDots9}},
	TileDots9:       {{TileDots7, TileDots8}},
	TileBamboo1:     {{TileBamboo2, TileBamboo3}},
	TileBamboo2:     {{TileBamboo1, TileBamboo3}, {TileBamboo3, TileBamboo4}},
	TileBamboo3:     {{TileBamboo1, TileBamboo2}, {TileBamboo2, TileBamboo4}, {TileBamboo4, TileBamboo5}},
	TileBamboo4:     {{TileBamboo2, TileBamboo3}, {TileBamboo3, TileBamboo5}, {TileBamboo5, TileBamboo6}},
	TileBamboo5:     {{TileBamboo3, TileBamboo4}, {TileBamboo4, TileBamboo6}, {TileBamboo6, TileBamboo7}},
	TileBamboo6:     {{TileBamboo4, TileBamboo5}, {TileBamboo5, TileBamboo7}, {TileBamboo7, TileBamboo8}},
	TileBamboo7:     {{TileBamboo5, TileBamboo6}, {TileBamboo6, TileBamboo8}, {TileBamboo8, TileBamboo9}},
	TileBamboo8:     {{TileBamboo6, TileBamboo7}, {TileBamboo7, TileBamboo9}},
	TileBamboo9:     {{TileBamboo7, TileBamboo8}},
	TileCharacters1: {{TileCharacters2, TileCharacters3}},
	TileCharacters2: {{TileCharacters1, TileCharacters3}, {TileCharacters3, TileCharacters4}},
	TileCharacters3: {{TileCharacters1, TileCharacters2}, {TileCharacters2, TileCharacters4}, {TileCharacters4, TileCharacters5}},
	TileCharacters4: {{TileCharacters2, TileCharacters3}, {TileCharacters3, TileCharacters5}, {TileCharacters5, TileCharacters6}},
	TileCharacters5: {{TileCharacters3, TileCharacters4}, {TileCharacters4, TileCharacters6}, {TileCharacters6, TileCharacters7}},
	TileCharacters6: {{TileCharacters4, TileCharacters5}, {TileCharacters5, TileCharacters7}, {TileCharacters7, TileCharacters8}},
	TileCharacters7: {{TileCharacters5, TileCharacters6}, {TileCharacters6, TileCharacters8}, {TileCharacters8, TileCharacters9}},
	TileCharacters8: {{TileCharacters6, TileCharacters7}, {TileCharacters7, TileCharacters9}},
	TileCharacters9: {{TileCharacters7, TileCharacters8}},
}

func isValidSequence(tile0, tile1, tile2 Tile) bool {
	others, ok := sequences[tile0]
	if !ok {
		return false
	}
	for _, tiles := range others {
		if (tile1 == tiles[0] && tile2 == tiles[1]) || (tile1 == tiles[1] && tile2 == tiles[0]) {
			return true
		}
	}
	return false
}
