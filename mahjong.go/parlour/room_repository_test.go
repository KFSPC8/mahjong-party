package parlour

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/mahjong.go"
)

var (
	conn *pgx.Conn
)

func getTx() pgx.Tx {
	if conn == nil {
		var err error
		conn, err = pgx.Connect(context.Background(), "postgres://localhost/parlour_test?sslmode=disable")
		if err != nil {
			panic(err)
		}
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	return tx
}

func TestPostgresRoomRepository(t *testing.T) {
	t.Run("saves and retrieves room", func(t *testing.T) {
		tx := getTx()
		defer tx.Rollback(context.Background())

		repo := NewPostgresRoomRepository(tx)
		room := &Room{
			ID: "ABCD",
			Results: []mahjong.Result{
				{
					Dealer:       1,
					Wind:         1,
					Winner:       1,
					Loser:        -1,
					Points:       1,
					WinningTiles: []mahjong.Tile{mahjong.TileDragonsWhite},
				},
			},
			clients: map[chan RoomView]string{},
		}
		err := repo.Save(room)
		assert.NoError(t, err)

		got, err := repo.Get("ABCD")
		assert.NoError(t, err)
		assert.Equal(t, room, got)
	})
	t.Run("sets room ID on save if unset", func(t *testing.T) {
		tx := getTx()
		defer tx.Rollback(context.Background())

		repo := NewPostgresRoomRepository(tx)
		room := NewRoom(Player{
			ID:   "iPRk13H8j/MHaP3vhCjnAg",
			Name: "Jiayu",
		})
		err := repo.Save(room)
		assert.NoError(t, err)

		got, err := repo.Get(room.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, got.ID)
	})
	t.Run("generates new ID on conflict", func(t *testing.T) {
		tx := getTx()
		defer tx.Rollback(context.Background())

		repo := NewPostgresRoomRepository(tx)
		room := NewRoom(Player{
			ID:   "iPRk13H8j/MHaP3vhCjnAg",
			Name: "Jiayu",
		})
		err := repo.Save(room)
		assert.NoError(t, err)
		assert.NotEmpty(t, room.ID)
	})
}
