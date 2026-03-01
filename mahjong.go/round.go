package mahjong

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

const (
	MinTilesLeft = 16
)

var (
	ErrNoMoreRounds = errors.New("no more rounds")
)

// Round represents a round in a mahjong game.
type Round struct {
	// Scores contains the score for each player in the game.
	Scores [4]int

	// Hands contains the corresponding hand for each player in the game.
	Hands [4]Hand

	// Wall contains the remaining tiles left to be drawn.
	Wall []Tile

	// Discards contains all the previously discarded tiles.
	Discards []Tile

	// Wind is the prevailing wind for the round.
	Wind Direction

	// Dealer is the integer offset of the dealer for the round.
	Dealer int

	// Turn is the integer offset of the player whose turn it currently is.
	Turn int

	// Phase is the current turn phase.
	Phase Phase

	// Events contains all the events that happened in the round.
	Events []Event

	// Result is the outcome of the round.
	Result *Result

	// Rules are the scoring rules for this round.
	Rules Rules

	// Finished indicates whether a round is over.
	Finished bool

	// WinningTile is used to determine if another player with higher
	// precedence can hu after someone else has already done so.
	WinningTile Tile

	LastActionTime   time.Time
	ReservedDuration time.Duration
}

func (r *Round) lastDiscard() Tile {
	if len(r.Discards) == 0 {
		return ""
	}
	return r.Discards[len(r.Discards)-1]
}

func (r *Round) popLastDiscard() Tile {
	tile := r.Discards[len(r.Discards)-1]
	r.Discards = r.Discards[:len(r.Discards)-1]
	return tile
}

func (r *Round) drawFront() Tile {
	drawn := r.Wall[0]
	r.Wall = r.Wall[1:]
	return drawn
}

func (r *Round) drawFrontN(n int) []Tile {
	draws := r.Wall[:n]
	r.Wall = r.Wall[n:]
	return draws
}

func (r *Round) drawBack() Tile {
	drawn := r.Wall[len(r.Wall)-1]
	r.Wall = r.Wall[:len(r.Wall)-1]
	return drawn
}

func (r *Round) previousTurn() int {
	return (r.Turn + 3) % 4
}

func (r *Round) replaceTile(seat int, t time.Time) {
	drawn := r.drawBack()
	for isFlower(drawn) {
		r.addFlower(seat, t, drawn)
		drawn = r.drawBack()
	}
	r.Hands[seat].Concealed.Add(drawn)
}

func (r *Round) seatWind(seat int) Direction {
	return Direction((seat - r.Dealer + 4) % 4)
}

func (r *Round) Draw(seat int, t time.Time) error {
	if r.Turn != seat {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDraw {
		return errors.New("wrong phase")
	}
	if t.Before(r.LastActionTime.Add(r.ReservedDuration)) {
		return errors.New("cannot draw during reserved duration")
	}
	r.Events = append(r.Events, Event{
		Type: EventDraw,
		Seat: seat,
		Time: timeInMillis(t),
	})
	drawn := r.drawFront()
	for isFlower(drawn) {
		r.addFlower(seat, t, drawn)
		drawn = r.drawBack()
	}
	hand := &r.Hands[seat]
	hand.Concealed.Add(drawn)
	r.Phase = PhaseDiscard
	r.LastActionTime = t
	return nil
}

func (r *Round) Discard(seat int, t time.Time, tile Tile) error {
	if seat != r.Turn {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDiscard {
		return errors.New("wrong phase")
	}
	if !r.Hands[seat].Concealed.Contains(tile) {
		return errors.New("missing tiles")
	}
	if len(r.Wall) <= MinTilesLeft-1 {
		return errors.New("no draws left")
	}
	r.Hands[seat].Concealed.Remove(tile)
	r.Discards = append(r.Discards, tile)
	r.Turn = (r.Turn + 1) % 4
	r.Phase = PhaseDraw
	r.Events = append(r.Events, newEvent(EventDiscard, seat, t, tile))
	r.LastActionTime = t
	return nil
}

func (r *Round) Chi(seat int, t time.Time, tile1, tile2 Tile) error {
	if r.Finished {
		return errors.New("round finished")
	}
	if r.Turn != seat {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDraw {
		return errors.New("wrong phase")
	}
	if len(r.Discards) == 0 {
		return errors.New("no discards")
	}
	tile0 := r.lastDiscard()
	if !isValidSequence(tile0, tile1, tile2) {
		return errors.New("invalid sequence")

	}
	hand := &r.Hands[seat]
	if !hand.Concealed.Contains(tile1) || !hand.Concealed.Contains(tile2) {
		return errors.New("missing tiles")
	}
	if t.Before(r.LastActionTime.Add(r.ReservedDuration)) {
		return errors.New("cannot chi during reserved duration")
	}
	hand.Concealed.Remove(tile1)
	hand.Concealed.Remove(tile2)
	r.popLastDiscard()
	seq := []Tile{tile0, tile1, tile2}
	sort.Slice(seq, func(i, j int) bool {
		return seq[i] < seq[j]
	})
	hand.Revealed = append(hand.Revealed, Meld{
		Type:  MeldChi,
		Tiles: seq,
	})
	r.Phase = PhaseDiscard
	r.Events = append(r.Events, newEvent(EventChi, seat, t, seq...))
	r.LastActionTime = t
	return nil
}

func (r *Round) Pong(seat int, t time.Time) error {
	if r.Finished {
		return errors.New("round finished")
	}
	if seat == r.previousTurn() {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDraw {
		return errors.New("wrong phase")
	}
	if len(r.Discards) == 0 {
		return errors.New("no discards")
	}
	hand := &r.Hands[seat]
	if hand.Concealed.Count(r.lastDiscard()) < 2 {
		return errors.New("missing tiles")
	}
	tile := r.popLastDiscard()
	hand.Concealed.RemoveN(tile, 2)
	hand.Revealed = append(hand.Revealed, Meld{
		Type:  MeldPong,
		Tiles: []Tile{tile},
	})
	r.Events = append(r.Events, newEvent(EventPong, seat, t, tile))
	r.Turn = seat
	r.Phase = PhaseDiscard
	r.LastActionTime = t
	return nil
}

func (r *Round) GangFromDiscard(seat int, t time.Time) error {
	if r.Finished {
		return errors.New("round finished")
	}
	if seat == r.previousTurn() {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDraw {
		return errors.New("wrong phase")
	}
	if len(r.Discards) == 0 {
		return errors.New("no discards")
	}
	hand := &r.Hands[seat]
	if hand.Concealed.Count(r.lastDiscard()) < 3 {
		return errors.New("missing tiles")
	}
	tile := r.popLastDiscard()
	hand.Concealed.RemoveN(tile, 3)
	hand.Revealed = append(hand.Revealed, Meld{
		Type:  MeldGang,
		Tiles: []Tile{tile},
	})
	r.replaceTile(seat, t)
	r.Events = append(r.Events, newEvent(EventGang, seat, t, tile))
	r.Turn = seat
	r.Phase = PhaseDiscard
	r.LastActionTime = t
	return nil
}

func (r *Round) GangFromHand(seat int, t time.Time, tile Tile) error {
	if r.Finished {
		return errors.New("round finished")
	}
	if seat != r.Turn {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDiscard {
		return errors.New("wrong phase")
	}
	hand := &r.Hands[seat]
	if hand.Concealed.Count(tile) > 3 {
		hand.Concealed.RemoveN(tile, 4)
		hand.Revealed = append(hand.Revealed, Meld{
			Type:  MeldGang,
			Tiles: []Tile{tile},
		})
		r.replaceTile(seat, t)
		r.Events = append(r.Events, newEvent(EventGang, seat, t, tile))
		r.LastActionTime = t
		return nil
	}
	for i, meld := range hand.Revealed {
		if meld.Type == MeldPong && meld.Tiles[0] == tile && hand.Concealed.Count(tile) > 0 {
			hand.Concealed.Remove(tile)
			hand.Revealed[i].Type = MeldGang
			r.replaceTile(seat, t)
			r.Events = append(r.Events, newEvent(EventGang, seat, t, tile))
			r.LastActionTime = t
			return nil
		}
	}
	return errors.New("missing tiles")
}

func bestHand(winningHands []Melds, round *Round, seat int) (Melds, int) {
	melds := append(round.Hands[seat].Revealed, winningHands[0]...)
	return winningHands[0], score(round, seat, melds)
}

func winningTiles(flowers []Tile, melds Melds, rest Melds) []Tile {
	melds = append(melds, rest...)
	sort.Sort(melds)
	return append(flowers, melds.Tiles()...)
}

func removeTile(tiles []Tile, tile Tile) []Tile {
	for i, t := range tiles {
		if t == tile {
			return append(tiles[:i], tiles[i+1:]...)
		}
	}
	return tiles
}

func (r *Round) tsumo(seat int) (best Melds, points int, err error) {
	if r.Finished {
		err = errors.New("already won")
		return
	}
	winningHands := search(r.Hands[seat].Concealed)
	if len(winningHands) == 0 {
		err = errors.New("missing tiles")
		return
	}
	best, points = bestHand(winningHands, r, seat)
	if points == 0 {
		err = errors.New("no tai")
		return
	}
	return
}

func (r *Round) ron(seat int, t time.Time) (best Melds, points, loser int, err error) {
	loser = r.previousTurn()
	if r.Finished {
		if t.After(r.LastActionTime.Add(r.ReservedDuration)) {
			err = errors.New("too late")
			return
		}
		winnerPrecedence := (r.Result.Winner - loser + 3) % 4
		precedence := (seat - loser + 3) % 4
		if precedence >= winnerPrecedence {
			err = errors.New("no precedence")
			return
		}
	}
	var winningTile Tile
	if r.WinningTile != "" {
		winningTile = r.WinningTile
	} else {
		winningTile = r.lastDiscard()
	}
	winningHands := search(r.Hands[seat].Concealed, winningTile)
	if len(winningHands) == 0 {
		err = errors.New("missing tiles")
		return
	}
	best, points = bestHand(winningHands, r, seat)
	if points == 0 {
		err = errors.New("no tai")
		return
	}
	if !r.Finished {
		// take the winning tile from the discard pile
		r.WinningTile = r.popLastDiscard()
	} else {
		// take it from the previous winner
		r.Hands[r.Result.Winner].Finished = removeTile(r.Hands[r.Result.Winner].Finished, r.WinningTile)
	}
	return
}

func (r *Round) Hu(seat int, t time.Time) error {
	if seat == r.previousTurn() {
		return errors.New("wrong turn")
	}
	if r.Turn != seat && r.Phase == PhaseDiscard {
		return errors.New("wrong turn")
	}
	var best Melds
	var points, loser int
	var err error
	if r.Phase == PhaseDiscard {
		loser = -1
		best, points, err = r.tsumo(seat)
	} else {
		best, points, loser, err = r.ron(seat, t)
	}
	if err != nil {
		return err
	}
	r.Hands[seat].Concealed = TileBag{}
	r.Hands[seat].Finished = best.Tiles()
	// undo previous score distribution if someone won previously
	if r.Result != nil {
		for i, delta := range winnings(r.Rules, r.Result.Winner, r.Result.Loser, r.Result.Points) {
			r.Scores[i] -= delta
		}
	}
	r.Result = &Result{
		Dealer:       r.Dealer,
		Wind:         r.Wind,
		Winner:       seat,
		WinningTiles: winningTiles(r.Hands[seat].Flowers, r.Hands[seat].Revealed, best),
		Loser:        loser,
		Points:       points,
	}
	r.LastActionTime = t
	r.Events = append(r.Events, newEvent(EventHu, seat, t))
	for i, delta := range winnings(r.Rules, seat, loser, points) {
		r.Scores[i] += delta
	}
	r.Finished = true
	return nil
}

func (r *Round) distributeTiles() {
	r.Hands = [4]Hand{
		{
			Flowers:   []Tile{},
			Revealed:  []Meld{},
			Concealed: TileBag{},
		},
		{
			Flowers:   []Tile{},
			Revealed:  []Meld{},
			Concealed: TileBag{},
		},
		{
			Flowers:   []Tile{},
			Revealed:  []Meld{},
			Concealed: TileBag{},
		},
		{
			Flowers:   []Tile{},
			Revealed:  []Meld{},
			Concealed: TileBag{},
		},
	}
	order := []int{r.Dealer, (r.Dealer + 1) % 4, (r.Dealer + 2) % 4, (r.Dealer + 3) % 4}
	// draw 4 tiles 3 times
	for i := 0; i < 3; i++ {
		var draws []Tile
		for _, seat := range order {
			draws = r.drawFrontN(4)
			r.Hands[seat].Concealed.Add(draws...)
		}
	}
	// draw one tile
	var draw Tile
	for _, seat := range order {
		draw = r.drawFront()
		r.Hands[seat].Concealed.Add(draw)
	}
	// dealer draws one extra tile
	draw = r.drawFront()
	r.Hands[r.Dealer].Concealed.Add(draw)
	// replace flowers
	for len(order) > 0 {
		seat := order[0]
		order = order[1:]
		mustReplaceAgain := false
		var flowers, replacements []Tile
		for tile := range r.Hands[seat].Concealed {
			if isFlower(tile) {
				draw = r.drawBack()
				if isFlower(draw) {
					mustReplaceAgain = true
				}
				flowers = append(flowers, tile)
				replacements = append(replacements, draw)
			}
		}
		if mustReplaceAgain {
			order = append(order, seat)
		}
		sort.Slice(flowers, func(i, j int) bool {
			return flowers[i] < flowers[j]
		})
		r.Hands[seat].Flowers = append(r.Hands[seat].Flowers, flowers...)
		r.Hands[seat].Concealed.Remove(flowers...)
		r.Hands[seat].Concealed.Add(replacements...)
	}
	return
}

func contains(tiles []Tile, tile Tile) bool {
	for _, t := range tiles {
		if t == tile {
			return true
		}
	}
	return false
}

func containsAll(haystack []Tile, needles []Tile) bool {
	for _, needle := range needles {
		if !contains(haystack, needle) {
			return false
		}
	}
	return true
}

func (r *Round) addFlower(seat int, t time.Time, flower Tile) {
	r.Events = append(r.Events, Event{
		Type:  EventFlower,
		Seat:  seat,
		Time:  timeInMillis(t),
		Tiles: []Tile{flower},
	})
	r.Hands[seat].Flowers = append(r.Hands[seat].Flowers, flower)
	if groups, ok := bites[flower]; ok {
		for _, group := range groups {
			if containsAll(r.Hands[seat].Flowers, group.Flowers) {
				r.Events = append(r.Events, Event{
					Type:  EventBitten,
					Seat:  seat,
					Time:  timeInMillis(t),
					Tiles: group.Flowers,
				})
				for i := range r.Scores {
					if i != seat {
						r.Scores[i] -= group.Payout
					}
				}
				r.Scores[seat] += 3 * group.Payout
			}
		}
	}
}

func (r *Round) Start(seed int64, t time.Time) {
	r.Wall = newWall(rand.New(rand.NewSource(seed)))
	r.distributeTiles()
	r.Turn = r.Dealer
	r.Phase = PhaseDiscard
	r.Discards = []Tile{}
	r.LastActionTime = t
	r.Events = []Event{newEvent(EventStart, 0, t)}
}

// Next returns a new round, setting the dealer and the prevailing wind
// depending on the outcome of this round.
func (r *Round) Next() (*Round, error) {
	if !r.Finished {
		return nil, errors.New("unfinished")
	}
	dealer := r.Dealer
	wind := r.Wind
	if r.Result.Winner != dealer {
		if dealer == 3 && wind == DirectionNorth {
			return nil, ErrNoMoreRounds
		}
		dealer = (r.Dealer + 1) % 4
		if dealer == 0 {
			wind++
		}
	}
	return &Round{
		Scores:           r.Scores,
		Dealer:           dealer,
		Wind:             wind,
		Rules:            r.Rules,
		ReservedDuration: r.ReservedDuration,
	}, nil
}

// End ends a round in a draw. Only the player who drew the last available tile
// from the wall may initiate this action.
func (r *Round) End(seat int, t time.Time) error {
	if r.Turn != seat {
		return errors.New("wrong turn")
	}
	if r.Phase != PhaseDiscard {
		return errors.New("wrong phase")
	}
	if len(r.Wall) >= MinTilesLeft {
		return errors.New("some draws remaining")
	}
	r.Finished = true
	r.Result = &Result{
		Dealer: r.Dealer,
		Wind:   r.Wind,
		Winner: -1,
		Loser:  -1,
	}
	r.LastActionTime = t
	r.Events = append(r.Events, newEvent(EventEnd, seat, t))
	return nil
}

// View returns a view of a round from a certain seat. Values of seat outside
// of [0, 3] will return a bystander's view of the round.
func (r *Round) View(seat int) RoundView {
	var hands [4]Hand
	for i, hand := range r.Hands {
		if seat == i {
			hands[i] = hand
		} else {
			hands[i] = hand.View()
		}
	}
	return RoundView{
		Seat:             seat,
		Scores:           r.Scores,
		Hands:            hands,
		DrawsLeft:        len(r.Wall) - MinTilesLeft + 1,
		Discards:         r.Discards,
		Wind:             r.Wind,
		Dealer:           r.Dealer,
		Turn:             r.Turn,
		Phase:            r.Phase,
		Events:           r.Events,
		Result:           r.Result,
		LastActionTime:   r.LastActionTime.UnixNano() / 1e6,
		ReservedDuration: r.ReservedDuration.Milliseconds(),
		Finished:         r.Finished,
	}
}

func newWall(r *rand.Rand) []Tile {
	var wall []Tile
	wall = append(wall, flowerTiles...)
	for _, tile := range suitedTiles {
		wall = append(wall, tile, tile, tile, tile)
	}
	r.Shuffle(len(wall), func(i, j int) {
		wall[i], wall[j] = wall[j], wall[i]
	})
	return wall
}
