package parlour

import (
	"fmt"
	"time"

	"github.com/yi-jiayu/mahjong.go"
)

type AI interface {
	Think(view RoomView) *Action
}

type discardRandomTileAI struct{}

func (ai discardRandomTileAI) Think(view RoomView) *Action {
	if view.Round == nil {
		return nil
	}
	if view.Round.Turn != view.Round.Seat {
		return nil
	}
	if view.Round.Phase == mahjong.PhaseDraw {
		time.Sleep(time.Duration(view.Round.ReservedDuration)*time.Millisecond + time.Second)
		return &Action{
			Nonce: view.Nonce,
			Type:  ActionDraw,
		}
	} else if view.Round.Phase == mahjong.PhaseDiscard {
		if view.Round.DrawsLeft <= 0 && !view.Round.Finished {
			return &Action{
				Nonce: view.Nonce,
				Type:  ActionEndRound,
			}
		}
		var tileToDiscard mahjong.Tile
		for tile := range view.Round.Hands[view.Round.Seat].Concealed {
			tileToDiscard = tile
			break
		}
		return &Action{
			Nonce: view.Nonce,
			Type:  ActionDiscard,
			Tiles: []mahjong.Tile{tileToDiscard},
		}
	}
	return nil
}

type Bot struct {
	ID      string
	Room    *Room
	Updates chan RoomView
	AI      AI
}

func (b *Bot) Start(roomService *roomService) {
	for view := range b.Updates {
		go func(view RoomView) {
			action := b.AI.Think(view)
			if action == nil {
				return
			}
			err := roomService.Dispatch(b.Room, b.ID, *action)
			if err != nil && err != errInvalidNonce {
				fmt.Printf("room=%s bot=%s error making move: %v", b.Room.ID, b.ID, err)
			}
		}(view)
	}
}
