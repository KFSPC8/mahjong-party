package mahjong

// RoundView represents a player's view of a round.
type RoundView struct {
	Seat      int       `json:"seat"`
	Scores    [4]int    `json:"scores"`
	Hands     [4]Hand   `json:"hands"`
	DrawsLeft int       `json:"draws_left"`
	Discards  []Tile    `json:"discards"`
	Wind      Direction `json:"wind"`
	Dealer    int       `json:"dealer"`
	Turn      int       `json:"turn"`
	Phase     Phase     `json:"phase"`
	Events    []Event   `json:"events"`
	Result    *Result   `json:"result,omitempty"`
	Finished  bool      `json:"finished"`

	// LastActionTime is the time the last action took place represented in milliseconds since the Unix epoch.
	LastActionTime int64 `json:"last_action_time"`

	// ReservedDuration is a duration in milliseconds reserved for players to pong or gang after a discard.
	ReservedDuration int64 `json:"reserved_duration"`
}
