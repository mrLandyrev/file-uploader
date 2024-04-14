package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/mrLandyrev/file-uploader/internal/usecases"
)

func Bootstrap(db *sql.DB) (*gin.Engine, error) {
	server := gin.Default()
	saveFileUseCase := usecases.BuildSaveFileUseCase(db)

	server.Use(BuildCORSMiddleware())

	server.GET("/email/list", BuildHandleList(db))
	server.POST("/email/upload", BuildHandlerUpload(saveFileUseCase))
	server.GET("/email/:id", BuildHandleDetails(db))

	return server, nil
}
