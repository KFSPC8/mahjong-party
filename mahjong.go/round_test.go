package mahjong

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRound_Draw(t *testing.T) {
	t.Run("cannot draw on wrong turn", func(t *testing.T) {
		r := &Round{Turn: 0}
		err := r.Draw(1, time.Now())
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only draw during draw phase", func(t *testing.T) {
		r := &Round{Turn: 0, Phase: PhaseDiscard}
		err := r.Draw(0, time.Now())
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("cannot draw during reserved duration", func(t *testing.T) {
		now := time.Now()
		oneSecondAgo := now.Add(-time.Second)
		r := &Round{
			Turn:             0,
			Phase:            PhaseDraw,
			LastActionTime:   oneSecondAgo,
			ReservedDuration: 2 * time.Second,
		}
		err := r.Draw(0, time.Now())
		assert.EqualError(t, err, "cannot draw during reserved duration")
	})
	t.Run("successful draw", func(t *testing.T) {
		seat := 0
		r := &Round{
			Wall:  []Tile{TileBamboo1, TileDots5},
			Turn:  seat,
			Phase: PhaseDraw,
			Hands: [4]Hand{{Concealed: NewTileBag([]Tile{TileWindsWest})}},
		}
		now := time.Now()
		err := r.Draw(seat, now)
		assert.NoError(t, err)
		assert.Equal(t, []Tile{TileDots5}, r.Wall)
		assert.Equal(t, NewTileBag([]Tile{TileWindsWest, TileBamboo1}), r.Hands[seat].Concealed)
		assert.Equal(t, seat, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Equal(t, []Event{{
			Type: EventDraw,
			Seat: seat,
			Time: timeInMillis(now),
		}}, r.Events)
		assert.Equal(t, now, r.LastActionTime)
	})
	t.Run("successful draw with flowers", func(t *testing.T) {
		seat := 0
		r := &Round{
			Wall:  []Tile{TileGentlemen1, TileBamboo1, TileDots5, TileGentlemen2},
			Turn:  seat,
			Phase: PhaseDraw,
			Hands: [4]Hand{{Concealed: NewTileBag([]Tile{TileWindsWest})}},
		}
		now := time.Now()
		err := r.Draw(seat, now)
		assert.NoError(t, err)
		assert.Equal(t, []Tile{TileGentlemen1, TileGentlemen2}, r.Hands[seat].Flowers)
		assert.Equal(t, []Tile{TileBamboo1}, r.Wall)
		assert.Equal(t, NewTileBag([]Tile{TileWindsWest, TileDots5}), r.Hands[seat].Concealed)
		assert.Equal(t, seat, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Contains(t, r.Events, Event{
			Type: EventDraw,
			Seat: seat,
			Time: timeInMillis(now),
		})
	})
}

func TestRound_Discard(t *testing.T) {
	t.Run("can only discard during own turn", func(t *testing.T) {
		r := &Round{Turn: 0}
		err := r.Discard(1, time.Now(), "")
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only discard during discard phase", func(t *testing.T) {
		r := &Round{Turn: 0, Phase: PhaseDraw}
		err := r.Discard(0, time.Now(), "")
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("cannot discard when missing tile", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{}},
		}
		err := r.Discard(0, time.Now(), TileDragonsRed)
		assert.EqualError(t, err, "missing tiles")
	})
	t.Run("cannot discard when no draws left", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{Concealed: TileBag{TileDragonsRed: 1}}},
		}
		err := r.Discard(0, time.Now(), TileDragonsRed)
		assert.EqualError(t, err, "no draws left")
	})
	t.Run("successful discard", func(t *testing.T) {
		seat := 0
		r := &Round{
			Wall: []Tile{
				"38八万", "35五万", "27六索", "44红中",
				"22一索", "34四万", "35五万", "20八筒",
				"37七万", "13一筒", "43北风", "26五索",
				"21九筒", "25四索", "42西风", "17五筒",
			},
			Turn:     seat,
			Phase:    PhaseDiscard,
			Discards: []Tile{TileWindsEast},
			Hands:    [4]Hand{{Concealed: NewTileBag([]Tile{TileCharacters1, TileWindsNorth})}},
		}
		now := time.Now()
		err := r.Discard(seat, now, TileWindsNorth)
		assert.NoError(t, err)
		assert.Equal(t, NewTileBag([]Tile{TileCharacters1}), r.Hands[seat].Concealed)
		assert.Equal(t, []Tile{TileWindsEast, TileWindsNorth}, r.Discards)
		assert.Equal(t, 1, r.Turn)
		assert.Equal(t, PhaseDraw, r.Phase)
		assert.Equal(t, []Event{{
			Type:  EventDiscard,
			Seat:  seat,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileWindsNorth},
		}}, r.Events)
		assert.Equal(t, now, r.LastActionTime)
	})
}

func TestRound_Chi(t *testing.T) {
	t.Run("cannot chi on wrong turn", func(t *testing.T) {
		r := &Round{Turn: 0}
		err := r.Chi(1, time.Now(), "", "")
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only chi during draw phase", func(t *testing.T) {
		r := &Round{Turn: 0, Phase: PhaseDiscard}
		err := r.Chi(0, time.Now(), "", "")
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("cannot chi when no discards", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDraw,
		}
		err := r.Chi(0, time.Now(), TileBamboo1, TileBamboo2)
		assert.EqualError(t, err, "no discards")
	})
	t.Run("cannot chi non-suited tile", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileDragonsRed},
		}
		err := r.Chi(0, time.Now(), "", "")
		assert.EqualError(t, err, "invalid sequence")
	})
	t.Run("cannot chi with invalid sequence", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileBamboo3},
		}
		err := r.Chi(0, time.Now(), TileCharacters2, TileCharacters4)
		assert.EqualError(t, err, "invalid sequence")
	})
	t.Run("cannot chi when missing tiles", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileBamboo3},
			Hands:    [4]Hand{{}},
		}
		err := r.Chi(0, time.Now(), TileBamboo2, TileBamboo4)
		assert.EqualError(t, err, "missing tiles")
	})
	t.Run("cannot chi during reserved duration", func(t *testing.T) {
		now := time.Now()
		oneSecondAgo := now.Add(-time.Second)
		r := &Round{
			Turn:             0,
			Phase:            PhaseDraw,
			Discards:         []Tile{TileBamboo4, TileBamboo3},
			Hands:            [4]Hand{{Concealed: NewTileBag([]Tile{TileWindsWest, TileBamboo1, TileBamboo2})}},
			LastActionTime:   oneSecondAgo,
			ReservedDuration: 2 * time.Second,
		}
		err := r.Chi(0, time.Now(), TileBamboo1, TileBamboo2)
		assert.EqualError(t, err, "cannot chi during reserved duration")
	})
	t.Run("cannot chi after round is finished", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Finished: true,
		}
		err := r.Chi(0, time.Now(), "", "")
		assert.EqualError(t, err, "round finished")
	})
	t.Run("successful chi", func(t *testing.T) {
		now := time.Now()
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileBamboo4, TileBamboo3},
			Hands:    [4]Hand{{Concealed: NewTileBag([]Tile{TileWindsWest, TileBamboo1, TileBamboo2})}},
		}
		err := r.Chi(0, now, TileBamboo1, TileBamboo2)
		assert.NoError(t, err)
		assert.Equal(t, []Tile{TileBamboo4}, r.Discards)
		assert.Equal(t, NewTileBag([]Tile{TileWindsWest}), r.Hands[0].Concealed)
		assert.Equal(t, Melds{{
			Type:  MeldChi,
			Tiles: []Tile{TileBamboo1, TileBamboo2, TileBamboo3},
		}}, r.Hands[0].Revealed)
		assert.Equal(t, 0, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Equal(t, []Event{{
			Type:  EventChi,
			Seat:  0,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileBamboo1, TileBamboo2, TileBamboo3},
		}}, r.Events)
		assert.Equal(t, now, r.LastActionTime)
	})
}

func TestRound_Pong(t *testing.T) {
	t.Run("cannot pong immediately after own turn", func(t *testing.T) {
		r := &Round{Turn: 1}
		err := r.Pong(0, time.Now())
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only pong during draw phase", func(t *testing.T) {
		r := &Round{Turn: 0, Phase: PhaseDiscard}
		err := r.Pong(0, time.Now())
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("cannot pong when no discards", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDraw,
		}
		err := r.Pong(0, time.Now())
		assert.EqualError(t, err, "no discards")
	})
	t.Run("cannot pong when missing tiles", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileDragonsRed},
			Hands:    [4]Hand{{}},
		}
		err := r.Pong(0, time.Now())
		assert.EqualError(t, err, "missing tiles")
	})
	t.Run("cannot pong when round finished", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Finished: true,
		}
		err := r.Pong(0, time.Now())
		assert.EqualError(t, err, "round finished")
	})
	t.Run("successful pong", func(t *testing.T) {
		seat := 1
		r := &Round{
			Turn:     3,
			Phase:    PhaseDraw,
			Discards: []Tile{TileDots1, TileDragonsRed},
			Hands:    [4]Hand{{}, {Concealed: NewTileBag([]Tile{TileWindsWest, TileDragonsRed, TileDragonsRed})}},
		}
		now := time.Now()
		err := r.Pong(seat, now)
		assert.NoError(t, err)
		assert.Equal(t, []Tile{TileDots1}, r.Discards)
		assert.Equal(t, NewTileBag([]Tile{TileWindsWest}), r.Hands[seat].Concealed)
		assert.Equal(t, Melds{{
			Type:  MeldPong,
			Tiles: []Tile{TileDragonsRed},
		}}, r.Hands[seat].Revealed)
		assert.Equal(t, seat, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Equal(t, []Event{{
			Type:  EventPong,
			Seat:  seat,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileDragonsRed},
		}}, r.Events)
		assert.Equal(t, now, r.LastActionTime)
	})
}

func TestRound_GangFromDiscard(t *testing.T) {
	t.Run("cannot gang from discard immediately after own turn", func(t *testing.T) {
		r := &Round{Turn: 1}
		err := r.GangFromDiscard(0, time.Now())
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only gang from discard during draw phase", func(t *testing.T) {
		r := &Round{Turn: 0, Phase: PhaseDiscard}
		err := r.GangFromDiscard(0, time.Now())
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("cannot gang from discard when no discards", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDraw,
		}
		err := r.GangFromDiscard(0, time.Now())
		assert.EqualError(t, err, "no discards")
	})
	t.Run("cannot gang from discard when not enough tiles", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileDragonsRed},
			Hands:    [4]Hand{{Concealed: TileBag{TileDragonsRed: 2}}},
		}
		err := r.GangFromDiscard(0, time.Now())
		assert.EqualError(t, err, "missing tiles")
	})
	t.Run("cannot gang from discard when round finished", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Finished: true,
		}
		err := r.GangFromDiscard(0, time.Now())
		assert.EqualError(t, err, "round finished")
	})
	t.Run("successful gang from discard", func(t *testing.T) {
		seat := 1
		r := &Round{
			Wall:     []Tile{TileCharacters4, TileCharacters6, TileGentlemen1},
			Turn:     3,
			Phase:    PhaseDraw,
			Discards: []Tile{TileDots1, TileDragonsRed},
			Hands: [4]Hand{{}, {
				Flowers:   []Tile{TileCat},
				Concealed: NewTileBag([]Tile{TileWindsWest, TileDragonsRed, TileDragonsRed, TileDragonsRed}),
			}},
		}
		now := time.Now()
		err := r.GangFromDiscard(seat, now)
		assert.NoError(t, err)
		assert.Equal(t, []Tile{TileDots1}, r.Discards)
		assert.Equal(t, []Tile{TileCat, TileGentlemen1}, r.Hands[seat].Flowers)
		assert.Equal(t, NewTileBag([]Tile{TileWindsWest, TileCharacters6}), r.Hands[seat].Concealed)
		assert.Equal(t, Melds{{
			Type:  MeldGang,
			Tiles: []Tile{TileDragonsRed},
		}}, r.Hands[seat].Revealed)
		assert.Equal(t, seat, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Contains(t, r.Events, Event{
			Type:  EventGang,
			Seat:  seat,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileDragonsRed},
		})
		assert.Equal(t, now, r.LastActionTime)
	})
}

func TestRound_GangFromHand(t *testing.T) {
	t.Run("cannot gang from hand on wrong turn", func(t *testing.T) {
		r := &Round{Turn: 1}
		err := r.GangFromHand(0, time.Now(), "")
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only gang from hand during discard phase", func(t *testing.T) {
		r := &Round{Turn: 0, Phase: PhaseDraw}
		err := r.GangFromHand(0, time.Now(), "")
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("cannot gang from hand when not enough tiles and no corresponding pong", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{}},
		}
		err := r.GangFromHand(0, time.Now(), TileDragonsRed)
		assert.EqualError(t, err, "missing tiles")
	})
	t.Run("cannot gang from hand when round finished", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDiscard,
			Finished: true,
		}
		err := r.GangFromHand(0, time.Now(), "")
		assert.EqualError(t, err, "round finished")
	})
	t.Run("successful concealed gang", func(t *testing.T) {
		seat := 0
		r := &Round{
			Wall:  []Tile{TileCharacters1, TileDots4, TileCat},
			Turn:  0,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{
				Flowers:   []Tile{TileSeasons1},
				Concealed: TileBag{TileDragonsRed: 4},
			}},
		}
		now := time.Now()
		err := r.GangFromHand(seat, now, TileDragonsRed)
		assert.NoError(t, err)
		assert.Equal(t, Melds{{
			Type:  MeldGang,
			Tiles: []Tile{TileDragonsRed},
		}}, r.Hands[seat].Revealed)
		assert.Equal(t, TileBag{TileDots4: 1}, r.Hands[seat].Concealed)
		assert.Equal(t, []Tile{TileCharacters1}, r.Wall)
		assert.Equal(t, seat, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Contains(t, r.Events, Event{
			Type:  EventGang,
			Seat:  seat,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileDragonsRed},
		})
		assert.Equal(t, now, r.LastActionTime)
	})
	t.Run("successful promote pong to gang", func(t *testing.T) {
		seat := 0
		r := &Round{
			Wall:  []Tile{TileCharacters1, TileDots4, TileCat},
			Turn:  0,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{
				Flowers: []Tile{TileSeasons1},
				Revealed: []Meld{{
					Type:  MeldPong,
					Tiles: []Tile{TileDragonsRed},
				}},
				Concealed: TileBag{TileDragonsRed: 1},
			}},
		}
		now := time.Now()
		err := r.GangFromHand(seat, now, TileDragonsRed)
		assert.NoError(t, err)
		assert.Equal(t, Melds{{
			Type:  MeldGang,
			Tiles: []Tile{TileDragonsRed},
		}}, r.Hands[seat].Revealed)
		assert.Equal(t, TileBag{TileDots4: 1}, r.Hands[seat].Concealed)
		assert.Equal(t, []Tile{TileCharacters1}, r.Wall)
		assert.Equal(t, seat, r.Turn)
		assert.Equal(t, PhaseDiscard, r.Phase)
		assert.Contains(t, r.Events, Event{
			Type:  EventGang,
			Seat:  seat,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileDragonsRed},
		})
		assert.Equal(t, now, r.LastActionTime)
	})
}

func TestRound_Hu(t *testing.T) {
	t.Run("cannot hu immediately after discarding", func(t *testing.T) {
		r := &Round{Turn: 1}
		err := r.Hu(0, time.Now())
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("cannot hu during other player's discard phase", func(t *testing.T) {
		r := &Round{Turn: 2, Phase: PhaseDiscard}
		err := r.Hu(0, time.Now())
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("cannot hu when winning hand cannot be found", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{Concealed: TileBag{}}},
		}
		err := r.Hu(0, time.Now())
		assert.EqualError(t, err, "missing tiles")
	})
	t.Run("cannot hu when no tai", func(t *testing.T) {
		seat := 1
		r := &Round{
			Turn:  seat,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{},
				{
					Flowers:  []Tile{},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileBamboo8, TileBamboo8, TileBamboo8,
						TileDragonsWhite, TileDragonsWhite,
					}),
				},
			},
		}
		now := time.Now()
		err := r.Hu(seat, now)
		assert.EqualError(t, err, "no tai")
	})
	t.Run("successful zi mo hu", func(t *testing.T) {
		seat := 1
		r := &Round{
			Turn:  seat,
			Phase: PhaseDiscard,
			Hands: [4]Hand{{},
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileCharacters8, TileCharacters8, TileCharacters8,
						TileDragonsWhite, TileDragonsWhite,
					}),
				},
			},
		}
		now := time.Now()
		err := r.Hu(seat, now)
		assert.NoError(t, err)
		assert.True(t, r.Finished)
		assert.Equal(t, Melds{
			{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}},
		}, r.Hands[seat].Revealed)
		assert.Equal(t, TileBag{}, r.Hands[seat].Concealed)
		assert.Equal(t, []Tile{
			TileBamboo6, TileBamboo7, TileBamboo8,
			TileCharacters8, TileCharacters8, TileCharacters8,
			TileWindsWest, TileWindsWest, TileWindsWest,
			TileDragonsWhite, TileDragonsWhite,
		}, r.Hands[seat].Finished)
		assert.Equal(t, &Result{
			Winner: seat,
			WinningTiles: []Tile{
				"05梅", "01猫",
				"15三筒", "16四筒", "17五筒",
				"27六索", "28七索", "29八索",
				"38八万", "38八万", "38八万",
				"42西风", "42西风", "42西风",
				"46白板", "46白板",
			},
			Loser:  -1,
			Points: 1,
		}, r.Result)
		assert.Equal(t, now, r.LastActionTime)
		assert.Equal(
			t,
			[]Event{{
				Type: EventHu,
				Seat: seat,
				Time: timeInMillis(now),
			}},
			r.Events,
		)
		assert.Equal(t, winnings(r.Rules, seat, -1, 1), r.Scores)
	})
	t.Run("successful hu from discards", func(t *testing.T) {
		seat := 2
		r := &Round{
			Turn:     0,
			Phase:    PhaseDraw,
			Discards: []Tile{TileDragonsRed, TileDragonsWhite},
			Hands: [4]Hand{{}, {},
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileCharacters8, TileCharacters8, TileCharacters8,
						TileDragonsWhite,
					}),
				},
			},
		}
		err := r.Hu(seat, time.Now())
		assert.NoError(t, err)
		assert.Equal(t, []Tile{TileDragonsRed}, r.Discards)
		assert.True(t, r.Finished)
		assert.Equal(t, &Result{
			Winner: seat,
			WinningTiles: []Tile{
				"05梅", "01猫",
				"15三筒", "16四筒", "17五筒",
				"27六索", "28七索", "29八索",
				"38八万", "38八万", "38八万",
				"42西风", "42西风", "42西风",
				"46白板", "46白板",
			},
			Loser:  3,
			Points: 2,
		}, r.Result)
	})
	t.Run("cannot hu again after huing", func(t *testing.T) {
		r := &Round{
			Turn:     0,
			Phase:    PhaseDiscard,
			Result:   &Result{Winner: 0},
			Finished: true,
		}
		err := r.Hu(0, time.Now())
		assert.EqualError(t, err, "already won")
	})
	t.Run("cannot override someone with higher precedence", func(t *testing.T) {
		r := &Round{
			Turn:             0,
			Phase:            PhaseDraw,
			Discards:         []Tile{TileDragonsRed, TileDragonsWhite},
			ReservedDuration: time.Second,
			Hands: [4]Hand{{},
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileCharacters8, TileCharacters8, TileCharacters8,
						TileDragonsWhite,
					}),
				},
			},
		}
		_ = r.Hu(1, time.Now())
		err := r.Hu(2, time.Now())
		assert.EqualError(t, err, "no precedence")
	})
	t.Run("can NOT be overridden by another player with higher precedence after reserved duration", func(t *testing.T) {
		seat := 1
		r := &Round{
			Turn:     0, // means seat 3 discarded
			Phase:    PhaseDraw,
			Discards: []Tile{TileDragonsRed, TileDragonsWhite},
			Hands: [4]Hand{
				{},
				// seat 1 can also hu on the same tile
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileCharacters8, TileCharacters8, TileCharacters8,
						TileDragonsWhite,
					}),
				},
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileCharacters8, TileCharacters8, TileCharacters8,
						TileDragonsWhite,
					}),
				},
			},
		}
		now := time.Now()
		oneSecondLater := now.Add(time.Second)
		// seat 2 hus first
		_ = r.Hu(2, now)

		// player 1 hus, but reserved duration is 0 so it's too late
		err := r.Hu(seat, oneSecondLater)
		assert.EqualError(t, err, "too late")
	})
	t.Run("can be overridden by another player with higher precedence within reserved duration", func(t *testing.T) {
		seat := 1
		r := &Round{
			Turn:             0, // means seat 3 discarded
			Phase:            PhaseDraw,
			Discards:         []Tile{TileDragonsRed, TileDragonsWhite},
			ReservedDuration: 2 * time.Second,
			Hands: [4]Hand{
				{},
				// seat 1 can also hu on the same tile
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileDragonsWhite, TileDragonsWhite, // needs a pong
						TileCharacters8, TileCharacters8,
					}),
				},
				{
					Flowers:  []Tile{TileGentlemen1, TileCat},
					Revealed: []Meld{{Type: MeldChi, Tiles: []Tile{TileDots3, TileDots4, TileDots5}}},
					Concealed: NewTileBag([]Tile{
						TileBamboo6, TileBamboo7, TileBamboo8,
						TileWindsWest, TileWindsWest, TileWindsWest,
						TileCharacters8, TileCharacters8, TileCharacters8,
						TileDragonsWhite, // needs eyes
					}),
				},
			},
		}
		now := time.Now()
		oneSecondLater := now.Add(time.Second)
		// seat 2 hus first
		_ = r.Hu(2, now)

		// player 1 hus
		err := r.Hu(seat, oneSecondLater)
		assert.NoError(t, err)
		assert.Equal(t, []Tile{
			TileBamboo6, TileBamboo7, TileBamboo8,
			TileCharacters8, TileCharacters8, TileCharacters8,
			TileWindsWest, TileWindsWest, TileWindsWest,
			TileDragonsWhite, // winning tile was removed
		}, r.Hands[2].Finished)
		assert.Equal(t, []Tile{TileDragonsRed}, r.Discards)
		assert.True(t, r.Finished)
		assert.Equal(t, &Result{
			Winner: seat,
			WinningTiles: []Tile{
				TileGentlemen1, TileCat,
				TileDots3, TileDots4, TileDots5,
				TileBamboo6, TileBamboo7, TileBamboo8,
				TileWindsWest, TileWindsWest, TileWindsWest,
				TileDragonsWhite, TileDragonsWhite, TileDragonsWhite,
				TileCharacters8, TileCharacters8,
			},
			Points: 2,
			Loser:  3,
		}, r.Result)
		assert.Equal(t, [4]int{-2, 8, -2, -4}, r.Scores)
	})
}

func TestRound_View(t *testing.T) {
	r := &Round{
		Scores:           [4]int{4, 2, 0, 1},
		Dealer:           1,
		Wind:             DirectionNorth,
		ReservedDuration: 2 * time.Second,
	}
	var ms int64 = 1598707747116
	now := time.Unix(ms/1000, (ms%1000)*1e6)
	r.Start(0, now)
	_ = r.Discard(1, now, TileBamboo1)
	t.Run("view from seat", func(t *testing.T) {
		seat := 1
		view := r.View(seat)
		assert.Equal(
			t,
			RoundView{
				Seat:   seat,
				Scores: r.Scores,
				Hands: [4]Hand{
					{Flowers: []Tile{"07菊"}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
					{Flowers: []Tile{}, Revealed: []Meld{}, Concealed: TileBag{"16四筒": 1, "27六索": 1, "29八索": 1, "34四万": 2, "35五万": 1, "36六万": 2, "38八万": 2, "43北风": 1, "44红中": 1, "46白板": 1}},
					{Flowers: []Tile{}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
					{Flowers: []Tile{"06兰", "12冬"}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
				},
				DrawsLeft: len(r.Wall) - 15,
				Discards:  r.Discards,
				Wind:      r.Wind,
				Dealer:    r.Dealer,
				Turn:      r.Turn,
				Phase:     r.Phase,
				Events: []Event{
					{
						Type: EventStart,
						Time: ms,
					},
					{
						Type:  EventDiscard,
						Time:  ms,
						Seat:  1,
						Tiles: []Tile{TileBamboo1},
					},
				},

				Result:           r.Result,
				LastActionTime:   ms,
				ReservedDuration: r.ReservedDuration.Milliseconds(),
			},
			view,
		)
	})
	t.Run("bystander's view", func(t *testing.T) {
		view := r.View(-1)
		assert.Equal(
			t,
			RoundView{
				Seat:   -1,
				Scores: r.Scores,
				Hands: [4]Hand{
					{Flowers: []Tile{"07菊"}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
					{Flowers: []Tile{}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
					{Flowers: []Tile{}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
					{Flowers: []Tile{"06兰", "12冬"}, Revealed: []Meld{}, Concealed: TileBag{"": 13}},
				},
				DrawsLeft: len(r.Wall) - 15,
				Discards:  r.Discards,
				Wind:      r.Wind,
				Dealer:    r.Dealer,
				Turn:      r.Turn,
				Phase:     r.Phase,
				Events: []Event{
					{
						Type: EventStart,
						Time: ms,
					},
					{
						Type:  EventDiscard,
						Time:  ms,
						Seat:  1,
						Tiles: []Tile{TileBamboo1},
					},
				},
				Result:           r.Result,
				LastActionTime:   ms,
				ReservedDuration: r.ReservedDuration.Milliseconds(),
			},
			view,
		)
	})
}

func Test_newWall(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	got := newWall(r)
	want := []Tile{"38八万", "35五万", "27六索", "44红中", "22一索", "34四万", "35五万", "20八筒", "37七万", "13一筒", "43北风", "26五索", "21九筒", "25四索", "42西风", "17五筒", "38八万", "36六万", "16四筒", "43北风", "20八筒", "22一索", "37七万", "25四索", "42西风", "30九索", "19七筒", "06兰", "27六索", "07菊", "40东风", "32二万", "29八索", "36六万", "34四万", "46白板", "32二万", "15三筒", "17五筒", "37七万", "42西风", "14二筒", "43北风", "20八筒", "28七索", "45青发", "17五筒", "36六万", "34四万", "14二筒", "12冬", "46白板", "22一索", "40东风", "37七万", "28七索", "29八索", "16四筒", "39九万", "13一筒", "24三索", "01猫", "27六索", "40东风", "41南风", "34四万", "24三索", "31一万", "31一万", "25四索", "13一筒", "26五索", "15三筒", "14二筒", "18六筒", "24三索", "11秋", "19七筒", "45青发", "41南风", "44红中", "39九万", "27六索", "26五索", "10夏", "15三筒", "21九筒", "36六万", "41南风", "33三万", "29八索", "23二索", "28七索", "04蜈蚣", "32二万", "38八万", "29八索", "05梅", "39九万", "21九筒", "46白板", "33三万", "09春", "32二万", "25四索", "30九索", "39九万", "23二索", "02老鼠", "24三索", "44红中", "28七索", "45青发", "18六筒", "31一万", "14二筒", "43北风", "13一筒", "45青发", "30九索", "18六筒", "22一索", "31一万", "16四筒", "17五筒", "26五索", "23二索", "21九筒", "35五万", "42西风", "03公鸡", "35五万", "18六筒", "30九索", "46白板", "38八万", "40东风", "19七筒", "15三筒", "41南风", "33三万", "16四筒", "20八筒", "23二索", "08竹", "33三万", "19七筒", "44红中"}
	assert.Equal(t, want, got)
}

func TestRound_distributeTiles(t *testing.T) {
	r := &Round{
		Dealer: 1,
		Wall: []Tile{
			"38八万", "35五万", "27六索", "44红中", // dealer (1) draws first
			"22一索", "34四万", "35五万", "20八筒",
			"37七万", "13一筒", "43北风", "26五索",
			"21九筒", "25四索", "42西风", "17五筒",

			"38八万", "36六万", "16四筒", "43北风",
			"20八筒", "22一索", "37七万", "25四索",
			"42西风", "30九索", "19七筒", "06兰", // 3 draws a flower
			"27六索", "07菊", "40东风", "32二万", // 0 draws a flower

			"29八索", "36六万", "34四万", "46白板",
			"32二万", "15三筒", "17五筒", "37七万",
			"42西风", "14二筒", "43北风", "20八筒",
			"28七索", "45青发", "17五筒", "36六万",

			"34四万",
			"14二筒",
			"12冬", // 3 draws another flower
			"46白板",

			"22一索",

			"40东风", "37七万", "28七索", "29八索", "16四筒", "39九万", "13一筒", "24三索", "44红中", "27六索", "40东风", "41南风", "34四万", "24三索", "31一万", "31一万", "25四索", "13一筒", "26五索", "15三筒", "14二筒", "18六筒", "24三索", "11秋", "19七筒", "45青发", "41南风", "44红中", "39九万", "27六索", "26五索", "10夏", "15三筒", "21九筒", "36六万", "41南风", "33三万", "29八索", "23二索", "28七索", "04蜈蚣", "32二万", "38八万", "29八索", "05梅", "39九万", "21九筒", "46白板", "33三万", "09春", "32二万", "25四索", "30九索", "39九万", "23二索", "02老鼠", "24三索", "44红中", "28七索", "45青发", "18六筒", "31一万", "14二筒", "43北风", "13一筒", "45青发", "30九索", "18六筒", "22一索", "31一万", "16四筒", "17五筒", "26五索", "23二索", "21九筒", "35五万", "42西风", "03公鸡", "35五万", "18六筒", "30九索", "46白板", "38八万", "40东风", "19七筒", "15三筒", "41南风", "33三万", "16四筒", "20八筒",

			"23二索",        // 3 replaces the fourth flower
			"08竹",         // 3 replaces the third flower and gets a fourth flower
			"33三万",        // 0 replaces one tile
			"19七筒", "01猫", // 3 replaces two tiles and gets a third flower
		},
	}
	r.distributeTiles()
	assert.Equal(t,
		NewTileBag([]Tile{"38八万", "35五万", "27六索", "44红中", "38八万", "36六万", "16四筒", "43北风", "29八索", "36六万", "34四万", "46白板", "34四万", "22一索"}),
		r.Hands[1].Concealed)
	assert.Equal(t,
		NewTileBag([]Tile{"22一索", "34四万", "35五万", "20八筒", "20八筒", "22一索", "37七万", "25四索", "32二万", "15三筒", "17五筒", "37七万", "14二筒"}),
		r.Hands[2].Concealed)
	assert.Equal(t,
		NewTileBag([]Tile{
			"37七万", "13一筒", "43北风", "26五索",
			"42西风", "30九索", "19七筒",
			"42西风", "14二筒", "43北风", "20八筒",
			"19七筒", "23二索", // replaced tiles
		}),
		r.Hands[3].Concealed)
	assert.ElementsMatch(t, []Tile{"06兰", "12冬", "01猫", "08竹"}, r.Hands[3].Flowers)
	assert.Equal(t,
		NewTileBag([]Tile{
			"21九筒", "25四索", "42西风", "17五筒",
			"27六索", "40东风", "32二万",
			"28七索", "45青发", "17五筒", "36六万",
			"46白板",
			"33三万", // replaced tile
		}),
		r.Hands[0].Concealed)
	assert.ElementsMatch(t, []Tile{"07菊"}, r.Hands[0].Flowers)
	assert.Equal(t,
		[]Tile{"40东风", "37七万", "28七索", "29八索", "16四筒", "39九万", "13一筒", "24三索", "44红中", "27六索", "40东风", "41南风", "34四万", "24三索", "31一万", "31一万", "25四索", "13一筒", "26五索", "15三筒", "14二筒", "18六筒", "24三索", "11秋", "19七筒", "45青发", "41南风", "44红中", "39九万", "27六索", "26五索", "10夏", "15三筒", "21九筒", "36六万", "41南风", "33三万", "29八索", "23二索", "28七索", "04蜈蚣", "32二万", "38八万", "29八索", "05梅", "39九万", "21九筒", "46白板", "33三万", "09春", "32二万", "25四索", "30九索", "39九万", "23二索", "02老鼠", "24三索", "44红中", "28七索", "45青发", "18六筒", "31一万", "14二筒", "43北风", "13一筒", "45青发", "30九索", "18六筒", "22一索", "31一万", "16四筒", "17五筒", "26五索", "23二索", "21九筒", "35五万", "42西风", "03公鸡", "35五万", "18六筒", "30九索", "46白板", "38八万", "40东风", "19七筒", "15三筒", "41南风", "33三万", "16四筒", "20八筒"},
		r.Wall)
}

func TestRound_Start(t *testing.T) {
	r := new(Round)
	now := time.Now()
	r.Start(0, now)
	assert.Equal(t, r.Dealer, r.Turn)
	assert.Equal(t, r.Phase, PhaseDiscard)
	assert.Equal(t, []Event{{Type: EventStart, Time: timeInMillis(now)}}, r.Events)
}

func TestRound_End(t *testing.T) {
	t.Run("can only end round during own turn", func(t *testing.T) {
		r := &Round{
			Turn: 1,
		}
		err := r.End(0, time.Now())
		assert.EqualError(t, err, "wrong turn")
	})
	t.Run("can only end round during discard phase", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDraw,
		}
		err := r.End(0, time.Now())
		assert.EqualError(t, err, "wrong phase")
	})
	t.Run("can only end round when there are less than 16 tiles in the wall", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDiscard,
			Wall: []Tile{
				"38八万", "35五万", "27六索", "44红中",
				"22一索", "34四万", "35五万", "20八筒",
				"37七万", "13一筒", "43北风", "26五索",
				"21九筒", "25四索", "42西风", "17五筒",
			},
		}
		err := r.End(0, time.Now())
		assert.EqualError(t, err, "some draws remaining")
	})
	t.Run("successfully ending round", func(t *testing.T) {
		r := &Round{
			Turn:  0,
			Phase: PhaseDiscard,
		}
		now := time.Now()
		err := r.End(0, now)
		assert.NoError(t, err)
		assert.True(t, r.Finished)
		assert.Equal(t, &Result{
			Dealer: r.Dealer,
			Wind:   r.Wind,
			Winner: -1,
			Loser:  -1,
		}, r.Result)
		assert.Equal(t, now, r.LastActionTime)
		assert.Equal(t, []Event{{Type: EventEnd, Seat: 0, Time: timeInMillis(now)}}, r.Events)
	})
}

func TestRound_Next(t *testing.T) {
	t.Run("cannot create next round when current round is not finished", func(t *testing.T) {
		r := &Round{Finished: false}
		_, err := r.Next()
		assert.EqualError(t, err, "unfinished")
	})
	t.Run("dealer does not win and dealer moves on", func(t *testing.T) {
		r := &Round{
			Finished: true,
			Dealer:   0,
			Result: &Result{
				Winner: 2,
			},
		}
		next, err := r.Next()
		assert.NoError(t, err)
		assert.Equal(t, 1, next.Dealer)
	})
	t.Run("dealer wins and remains dealer", func(t *testing.T) {
		r := &Round{
			Finished: true,
			Dealer:   2,
			Result: &Result{
				Winner: 2,
			},
		}
		next, err := r.Next()
		assert.NoError(t, err)
		assert.Equal(t, 2, next.Dealer)
	})
	t.Run("dealer moves on and prevailing wind changes", func(t *testing.T) {
		r := &Round{
			Finished: true,
			Dealer:   3,
			Wind:     DirectionEast,
			Result: &Result{
				Winner: 0,
			},
		}
		next, err := r.Next()
		assert.NoError(t, err)
		assert.Equal(t, 0, next.Dealer)
		assert.Equal(t, DirectionSouth, next.Wind)
	})
	t.Run("copies over round settings", func(t *testing.T) {
		r := &Round{
			Scores:           [4]int{4, 2, 1, -2},
			Finished:         true,
			ReservedDuration: 2 * time.Second,
			Rules:            RulesShooter,
			Result:           &Result{},
		}
		next, err := r.Next()
		assert.NoError(t, err)
		assert.Equal(t, r.Scores, next.Scores)
		assert.Equal(t, r.ReservedDuration, next.ReservedDuration)
		assert.Equal(t, r.Rules, next.Rules)
	})
	t.Run("no more rounds", func(t *testing.T) {
		r := &Round{
			Finished: true,
			Dealer:   3,
			Wind:     DirectionNorth,
			Result: &Result{
				Winner: 0,
			},
		}
		_, err := r.Next()
		assert.EqualError(t, err, "no more rounds")
	})
}

func TestRound_addFlowers(t *testing.T) {
	t.Run("adds flower to hand", func(t *testing.T) {
		r := &Round{
			Hands: [4]Hand{{Flowers: []Tile{TileGentlemen1}}},
		}
		now := time.Now()
		r.addFlower(0, now, TileCat)
		assert.Equal(t, []Tile{TileGentlemen1, TileCat}, r.Hands[0].Flowers)
		assert.Equal(t, []Event{{
			Type:  EventFlower,
			Seat:  0,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileCat},
		}}, r.Events)
	})
	t.Run("triggers bitten event", func(t *testing.T) {
		r := &Round{
			Hands: [4]Hand{{Flowers: []Tile{TileRat}}},
		}
		now := time.Now()
		r.addFlower(0, now, TileCat)
		assert.Contains(t, r.Events, Event{
			Type:  EventBitten,
			Seat:  0,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileCat, TileRat},
		})
		assert.Equal(t, [4]int{6, -2, -2, -2}, r.Scores)
	})
	t.Run("cumulative bitten events", func(t *testing.T) {
		r := &Round{
			Hands: [4]Hand{{Flowers: []Tile{TileCat, TileRat, TileRooster}}},
		}
		now := time.Now()
		r.addFlower(0, now, TileCentipede)
		assert.Contains(t, r.Events, Event{
			Type:  EventBitten,
			Seat:  0,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileRooster, TileCentipede},
		})
		assert.Contains(t, r.Events, Event{
			Type:  EventBitten,
			Seat:  0,
			Time:  timeInMillis(now),
			Tiles: []Tile{TileCat, TileRat, TileRooster, TileCentipede},
		})
		assert.Equal(t, [4]int{18, -6, -6, -6}, r.Scores)
	})
}
