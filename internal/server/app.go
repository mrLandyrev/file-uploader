package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/mrLandyrev/file-uploader/internal/database"
	"github.com/mrLandyrev/file-uploader/internal/handlers"
)

type App struct {
	db     *sql.DB
	engine *gin.Engine
}

func NewApp(config Config) (*App, error) {
	app := &App{}
	var err error

	if app.db, err = database.Bootstrap(); err != nil {
		return nil, err
	}

	if app.engine, err = handlers.Bootstrap(app.db); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() {
	a.engine.Run()
}
