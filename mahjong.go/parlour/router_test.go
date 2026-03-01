package parlour

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate ../bin/mockgen -destination mocks_test.go -package parlour -self_package github.com/yi-jiayu/mahjong.go/parlour . RoomRepository

func TestParlour_createRoomHandler(t *testing.T) {
	roomID := "ABCD"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	roomRepository := NewMockRoomRepository(ctrl)
	roomRepository.EXPECT().Save(gomock.Any()).DoAndReturn(func(room *Room) error {
		room.ID = roomID
		return nil
	})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	parlour := New(roomRepository, memstore.NewStore())
	parlour.configure(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/rooms", strings.NewReader("name=alice"))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, roomID, w.Body.String())
}

func TestParlour_joinRoomHandler(t *testing.T) {
	room := NewRoom(Player{Name: "alice"})
	room.ID = "ABCD"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	roomRepository := NewMockRoomRepository(ctrl)
	roomRepository.EXPECT().Get(room.ID).Return(room, nil)
	roomRepository.EXPECT().Save(gomock.Any()).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	parlour := New(roomRepository, memstore.NewStore())
	parlour.configure(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/rooms/%s/players", room.ID), strings.NewReader("name=bob"))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "bob", room.Players[1].Name)
}
