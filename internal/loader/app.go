package loader

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/mrLandyrev/file-uploader/internal/database"
	"github.com/mrLandyrev/file-uploader/internal/usecases"
	"golang.org/x/sync/semaphore"
)

type App struct {
	db   *sql.DB
	path string
}

func NewApp(config Config) (*App, error) {
	app := &App{}
	var err error

	if app.db, err = database.Bootstrap(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() {
	entries, err := os.ReadDir(a.path)
	if err != nil {
		log.Fatal(err)
	}

	steps := len(entries)

	maxWorkers := runtime.GOMAXPROCS(0)
	sem := semaphore.NewWeighted(int64(maxWorkers))
	ctx := context.TODO()

	for i, e := range entries {
		sem.Acquire(ctx, 1)
		go func(i int, n string) {
			defer sem.Release(1)
			f, _ := os.ReadFile(a.path + n)
			saveFileUseCase := usecases.BuildSaveFileUseCase(a.db)
			saveFileUseCase(usecases.ParseFiles(bytes.NewReader(f)), n)
			fmt.Printf("%d/%d\n", i, steps)
		}(i, e.Name())
	}
	sem.Acquire(ctx, int64(maxWorkers))
}
