package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/yi-jiayu/mahjong.go/parlour"
)

var (
	host, port string
	database   string
)

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "host to listen on")
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&database, "database", "", "database url")

	rand.Seed(time.Now().UnixNano())
}

func getKey(name string) []byte {
	authKey, err := base64.StdEncoding.DecodeString(os.Getenv(name))
	if err != nil {
		fmt.Printf("error decoding %s: %v\n", name, err)
		os.Exit(1)
	}
	if len(authKey) < 32 {
		fmt.Printf("error: %s should decode to 32 bytes of data\n", name)
		os.Exit(1)
	}
	return authKey
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL environment variable is required")
		os.Exit(1)
	}

	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("error connecting to postgres: %v\n", err)
		os.Exit(1)
	}

	roomRepository := parlour.NewPostgresRoomRepository(pool)
	authKey := getKey("PARLOUR_SESSION_AUTH_KEY")
	encKey := getKey("PARLOUR_SESSION_ENC_KEY")
	store := cookie.NewStore(authKey, encKey)
	p := parlour.New(roomRepository, store)
	err = p.Run(host + ":" + port)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
