package network

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"sync"
)

var (
	Pool     *pgxpool.Pool
	poolOnce sync.Once
)

func DbConnect() {
	poolOnce.Do(func() {
		var err error

		Pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

		if err != nil {
			log.Fatal("Unable to create pool database")
		}

		log.Println("Created pool database")

	})
}
