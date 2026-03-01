package parlour

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var errNotFound = errors.New("not found")

type RoomRepository interface {
	Save(room *Room) error
	Get(id string) (*Room, error)
}

func newRoomID() string {
	charset := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return string([]byte{
		charset[rand.Intn(len(charset))],
		charset[rand.Intn(len(charset))],
		charset[rand.Intn(len(charset))],
		charset[rand.Intn(len(charset))],
	})
}

type InMemoryRoomRepository struct {
	sync.RWMutex
	rooms map[string]*Room
}

func NewInMemoryRoomRepository() *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		rooms: map[string]*Room{},
	}
}

func (r *InMemoryRoomRepository) Save(room *Room) error {
	r.Lock()
	defer r.Unlock()
	for room.ID == "" {
		id := newRoomID()
		if _, exists := r.rooms[id]; exists {
			continue
		}
		room.ID = id
	}
	r.rooms[room.ID] = room
	return nil
}

func (r *InMemoryRoomRepository) Get(id string) (*Room, error) {
	r.RLock()
	defer r.RUnlock()
	id = strings.ToUpper(id)
	room, ok := r.rooms[id]
	if !ok {
		return nil, errNotFound
	}
	return room, nil
}

type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type PostgresRoomRepository struct {
	conn Conn
}

func (p *PostgresRoomRepository) Save(room *Room) error {
	ctx := context.Background()
	if room.ID == "" {
		for {
			id := newRoomID()
			tx, err := p.conn.Begin(ctx)
			if err != nil {
				return fmt.Errorf("error inserting room: %w", err)
			}
			_, err = tx.Exec(ctx, `insert into rooms (id, nonce, phase, players, round, results)
values ($1, $2, $3, $4, $5, $6)`,
				id,
				room.Nonce,
				room.Phase,
				room.Players,
				room.Round,
				room.Results,
			)
			if err != nil {
				var pgError *pgconn.PgError
				if errors.As(err, &pgError) {
					if pgError.Code == pgerrcode.UniqueViolation {
						_ = tx.Rollback(ctx)
						continue
					}
				}
				return fmt.Errorf("error inserting room: %w", err)
			}
			err = tx.Commit(ctx)
			if err != nil {
				return fmt.Errorf("error inserting room: %w", err)
			}
			room.ID = id
			return nil
		}
	}
	_, err := p.conn.Exec(ctx, `insert into rooms (id, nonce, phase, players, round, results)
values ($1, $2, $3, $4, $5, $6)
on conflict (id) do update set nonce=excluded.nonce,
                               phase=excluded.phase,
                               players=excluded.players,
                               round=excluded.round,
                               results=excluded.results`,
		room.ID,
		room.Nonce,
		room.Phase,
		room.Players,
		room.Round,
		room.Results,
	)
	if err != nil {
		return fmt.Errorf("error saving room: %w", err)
	}
	return nil
}

func (p *PostgresRoomRepository) Get(id string) (*Room, error) {
	var room Room
	err := p.conn.QueryRow(
		context.Background(),
		"select id, nonce, phase, players, round, results from rooms where id = $1", id,
	).Scan(&room.ID, &room.Nonce, &room.Phase, &room.Players, &room.Round, &room.Results)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting room: %w", err)
	}
	room.clients = make(map[chan RoomView]string)
	return &room, nil
}

func NewPostgresRoomRepository(conn Conn) *PostgresRoomRepository {
	return &PostgresRoomRepository{
		conn: conn,
	}
}
