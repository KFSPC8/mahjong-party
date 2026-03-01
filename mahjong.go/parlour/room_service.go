package parlour

import (
	"strings"
	"sync"
)

type Error struct {
	error
	internal bool
}

func (e *Error) Unwrap() error {
	return e.error
}

type roomService struct {
	RoomRepository RoomRepository

	cache map[string]*Room
	sync.Mutex
}

func (s *roomService) Get(id string) (*Room, error) {
	s.Lock()
	defer s.Unlock()
	id = strings.ToUpper(id)
	if room, ok := s.cache[id]; ok {
		return room, nil
	}
	room, err := s.RoomRepository.Get(id)
	if err != nil {
		return nil, &Error{error: err, internal: true}
	}
	s.cache[room.ID] = room

	// start bots
	for _, player := range room.Players {
		if player.IsBot {
			bot := Bot{
				ID:      player.ID,
				Room:    room,
				Updates: make(chan RoomView, 1),
				AI:      discardRandomTileAI{},
			}
			room.clients[bot.Updates] = bot.ID
			go bot.Start(s)
		}
	}
	room.broadcast()

	return room, nil
}

func (s *roomService) Create(host Player) (*Room, error) {
	s.Lock()
	defer s.Unlock()
	room := NewRoom(host)
	err := s.RoomRepository.Save(room)
	if err != nil {
		return nil, &Error{
			error:    err,
			internal: true,
		}
	}
	s.cache[room.ID] = room
	return room, nil
}

func (s *roomService) AddPlayer(room *Room, player Player) error {
	var svcErr error
	room.WithLock(func(r *Room) {
		err := room.addPlayer(player)
		if err != nil {
			svcErr = &Error{error: err}
			return
		}
		svcErr = s.RoomRepository.Save(r)
	})
	return svcErr
}

func (s *roomService) RemovePlayer(room *Room, playerID string) error {
	var svcErr error
	room.WithLock(func(r *Room) {
		room.removePlayer(playerID)
		svcErr = s.RoomRepository.Save(r)
	})
	return svcErr
}

func (s *roomService) Dispatch(room *Room, playerID string, action Action) error {
	var svcErr error
	room.WithLock(func(r *Room) {
		err := room.reduce(playerID, action)
		if err != nil {
			svcErr = &Error{error: err}
			return
		}
		svcErr = s.RoomRepository.Save(r)
	})
	return svcErr
}

var botNames = []string{"Francisco Bot", "Lupe Bot", "Mordecai Bot"}

func (s *roomService) AddBot(room *Room, playerID string) error {
	var svcErr error
	room.WithLock(func(r *Room) {
		if r.seat(playerID) == -1 {
			svcErr = &Error{error: errNotInRoom}
			return
		}
		if len(r.Players) > 3 {
			svcErr = &Error{error: errRoomFull}
			return
		}
		name := botNames[len(r.Players)-1]
		r.Players = append(r.Players, Player{
			ID:    name,
			Name:  name,
			IsBot: true,
		})
		r.broadcast()
		bot := Bot{
			ID:      name,
			Room:    r,
			Updates: make(chan RoomView),
			AI:      discardRandomTileAI{},
		}
		r.clients[bot.Updates] = bot.ID
		go bot.Start(s)
		svcErr = s.RoomRepository.Save(r)
	})
	return svcErr
}

func newRoomService(roomRepository RoomRepository) *roomService {
	return &roomService{
		RoomRepository: roomRepository,
		cache:          make(map[string]*Room),
	}
}
