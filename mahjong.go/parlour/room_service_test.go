package parlour

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_roomService_Get(t *testing.T) {
	t.Run("normalises room ID to upper case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		roomRepository := NewMockRoomRepository(ctrl)
		roomRepository.EXPECT().Get("ABCD").Return(&Room{ID: "ABCD"}, nil)

		service := newRoomService(roomRepository)
		_, _ = service.Get("abcd")
	})
	t.Run("caches room in memory", func(t *testing.T) {
		roomID := "ABCD"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		roomRepository := NewMockRoomRepository(ctrl)
		roomRepository.EXPECT().Get(gomock.Any()).Return(&Room{ID: roomID}, nil)

		service := newRoomService(roomRepository)
		room, err := service.Get(roomID)
		assert.NoError(t, err)

		got, err := service.Get(roomID)
		assert.NoError(t, err)
		assert.Same(t, room, got)
	})
}
