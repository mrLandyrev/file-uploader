package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/mrLandyrev/file-uploader/internal/usecases"
	"golang.org/x/sync/semaphore"
)

func main() {
	entries, err := os.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"192.168.0.105", 5432, "postgres", "aezakmi007", "postgres"))

	if err != nil {
		log.Fatal("Database connection error")
	}

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS email (
		id UUID not null,
		filename text not null,
		send_date timestamp not null,
		upload_date timestamp not null,
		"to" text not null,
		cc text not null,
		subject text not null,
		content text not null,
		"from" text not null,
		primary key(id)
	)`)

	if err != nil {
		log.Fatal("Database error", err)
	}

	steps := len(entries)

	maxWorkers := runtime.GOMAXPROCS(0)
	sem := semaphore.NewWeighted(int64(maxWorkers))
	ctx := context.TODO()

	for i, e := range entries {
		sem.Acquire(ctx, 1)
		go func(i int, n string) {
			defer sem.Release(1)
			f, _ := os.ReadFile(os.Args[1] + n)
			saveFileUseCase := usecases.BuildSaveFileUseCase(db)
			saveFileUseCase(usecases.ParseVor(bytes.NewReader(f)), n)
			fmt.Printf("%d/%d\n", i, steps)
		}(i, e.Name())
	}
	sem.Acquire(ctx, int64(maxWorkers))
}
