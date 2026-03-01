package mahjong

import (
	"time"
)

func timeInMillis(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// EventType represents the type of an event.
type EventType string

// Possible event types.
const (
	EventStart   = "start"
	EventDraw    = "draw"
	EventDiscard = "discard"
	EventChi     = "chi"
	EventPong    = "pong"
	EventGang    = "gang"
	EventHu      = "hu"
	EventEnd     = "end"
	EventFlower  = "flower"
	EventBitten  = "bitten"
)

// Event represents a player's view of an event.
type Event struct {
	// Type is the type of an event.
	Type EventType `json:"type"`

	// Seat is the integer offset of the player an event pertains to.
	Seat int `json:"seat"`

	// Time is the time an event occurred.
	Time int64 `json:"time"`

	// Tiles are the tiles involved in an event.
	Tiles []Tile `json:"tiles"`
}

func newEvent(eventType EventType, seat int, t time.Time, tiles ...Tile) Event {
	return Event{
		Type:  eventType,
		Seat:  seat,
		Time:  timeInMillis(t),
		Tiles: tiles,
	}
}
