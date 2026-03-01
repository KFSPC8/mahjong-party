package parlour

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yi-jiayu/mahjong.go"
)

const (
	KeySessionName = "session"
	KeyPlayerID    = "playerID"
	KeyRoom        = "room"
)

var sessionOptions = sessions.Options{
	Path:     "/",
	SameSite: http.SameSiteStrictMode,
	MaxAge:   2592000, // 30 days in seconds
}

// newPlayerID returns opaque string containing n bytes of entropy.
func newPlayerID(n int) string {
	data := make([]byte, n)
	_, err := mathRand.Read(data)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

func setPlayerID(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessionOptions)
	var playerID string
	v := session.Get(KeyPlayerID)
	if v == nil {
		playerID = newPlayerID(16)
	} else if p, ok := v.(string); !ok {
		playerID = newPlayerID(16)
	} else {
		playerID = p
	}
	c.Set(KeyPlayerID, playerID)
	session.Set(KeyPlayerID, playerID)
	err := session.Save()
	if err != nil {
		fmt.Printf("error saving session: %v\n", err)
	}
	c.Next()
}

func handleErrors(c *gin.Context) {
	c.Next()
	err := c.Errors.Last()
	if err == nil {
		return
	}
	var e *Error
	if errors.As(err.Err, &e) {
		if e.internal {
			fmt.Printf("internal error: %v", e)
			c.String(http.StatusInternalServerError, "internal error")
			return
		}
		_ = err.SetType(gin.ErrorTypePublic)
	}
	c.String(http.StatusBadRequest, err.Error())
}

func getName(c *gin.Context) (string, error) {
	name := c.PostForm("name")
	if name == "" {
		return "", errors.New("name is required")
	}
	if ok, _ := regexp.MatchString("[0-9A-Za-z ]+", name); !ok {
		return "", errors.New("name is invalid")
	}
	return name, nil
}

func (p *Parlour) createRoomHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerID := c.GetString(KeyPlayerID)
		name, err := getName(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		player := Player{
			ID:   playerID,
			Name: name,
		}
		room, err := p.roomService.Create(player)
		if err != nil {
			_ = c.Error(err)
			return
		}
		metricOpenRooms.Add(1)
		c.String(http.StatusCreated, room.ID)
	}
}

func (p *Parlour) joinRoomHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerID := c.GetString(KeyPlayerID)
		room := c.MustGet(KeyRoom).(*Room)
		name, err := getName(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		player := Player{
			ID:   playerID,
			Name: name,
		}
		err = p.roomService.AddPlayer(room, player)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func (p *Parlour) leaveRoomHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerID := c.GetString(KeyPlayerID)
		room := c.MustGet(KeyRoom).(*Room)
		err := p.roomService.RemovePlayer(room, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
}

func (p *Parlour) subscribeRoomHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerID := c.GetString(KeyPlayerID)
		room := c.MustGet(KeyRoom).(*Room)
		ch := make(chan RoomView, 1)
		room.AddClient(playerID, ch)
		metricRoomSubscriptions.Add(1)

		notify := c.Request.Context().Done()
		go func() {
			<-notify
			room.RemoveClient(ch)
			metricRoomSubscriptions.Add(-1)
		}()

		c.Stream(func(w io.Writer) bool {
			if view, ok := <-ch; ok {
				c.SSEvent("", view)
				return true
			}
			return false
		})
	}
}

func (p *Parlour) roomActionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerID := c.GetString(KeyPlayerID)
		room := c.MustGet(KeyRoom).(*Room)
		var action Action
		err := c.ShouldBindJSON(&action)
		if err != nil {
			_ = c.Error(err)
			return
		}
		err = p.roomService.Dispatch(room, playerID, action)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
}

func setConcealedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		room := c.MustGet(KeyRoom).(*Room)
		if room.Round == nil {
			c.String(http.StatusBadRequest, "round not started")
			return
		}
		seat, err := strconv.Atoi(c.Param("seat"))
		if err != nil {
			c.String(http.StatusBadRequest, "invalid seat")
			return
		}
		if seat < 0 || 3 < seat {
			c.String(http.StatusBadRequest, "invalid seat")
			return
		}
		var tiles mahjong.TileBag
		err = c.ShouldBindJSON(&tiles)
		if err != nil {
			_ = c.Error(err)
			return
		}
		room.Round.Hands[seat].Concealed = tiles
	}
}

func prependWallHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		room := c.MustGet(KeyRoom).(*Room)
		if room.Round == nil {
			c.String(http.StatusBadRequest, "round not started")
			return
		}
		tile := c.PostForm("tile")
		if tile == "" {
			c.String(http.StatusBadRequest, "tile is required")
			return
		}
		room.Round.Wall = append([]mahjong.Tile{mahjong.Tile(tile)}, room.Round.Wall...)
	}
}

func reshuffleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		room := c.MustGet(KeyRoom).(*Room)
		if room.Round == nil {
			c.String(http.StatusBadRequest, "round not started")
			return
		}
		room.Round = &mahjong.Round{
			Rules: mahjong.RulesDefault,
		}
		room.Round.Start(mathRand.Int63(), time.Now())
		room.broadcast()
	}
}

func (p *Parlour) addBotHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerID := c.GetString(KeyPlayerID)
		room := c.MustGet(KeyRoom).(*Room)
		err := p.roomService.AddBot(room, playerID)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
}

func (p *Parlour) setRoomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("roomID")
		room, err := p.roomService.Get(roomID)
		if errors.Is(err, errNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Set("room", room)
		c.Next()
	}
}

func (p Parlour) configure(r *gin.Engine) {
	r.Use(sessions.Sessions(KeySessionName, p.SessionStore))
	r.Use(setPlayerID)
	r.Use(handleErrors)
	r.POST("/rooms", p.createRoomHandler())
	room := r.Group("/rooms/:roomID")
	room.Use(p.setRoomMiddleware())
	{
		room.POST("/players", p.joinRoomHandler())
		room.DELETE("/players", p.leaveRoomHandler())
		room.GET("/live", p.subscribeRoomHandler())
		room.POST("/actions", p.roomActionsHandler())
		room.POST("/bots", p.addBotHandler())
		if gin.IsDebugging() {
			room.PUT("/round/hands/:seat/concealed", setConcealedHandler())
			room.POST("/round/wall", prependWallHandler())
			room.POST("/round/reshuffle", reshuffleHandler())
		}
	}
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
