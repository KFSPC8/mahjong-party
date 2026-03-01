package parlour

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Parlour struct {
	RoomRepository RoomRepository
	SessionStore   sessions.Store

	roomService *roomService
}

func New(roomRepository RoomRepository, sessionStore sessions.Store) *Parlour {
	p := &Parlour{
		RoomRepository: roomRepository,
		SessionStore:   sessionStore,
	}
	p.roomService = newRoomService(p.RoomRepository)
	return p
}

func (p *Parlour) Run(addr string) error {
	r := gin.Default()
	p.configure(r)
	return r.Run(addr)
}
