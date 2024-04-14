package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mrLandyrev/file-uploader/internal/handlers"
	"github.com/mrLandyrev/file-uploader/internal/usecases"
)

func main() {
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

	server := gin.Default()

	server.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	saveFileUseCase := usecases.BuildSaveFileUseCase(db)

	server.GET("/email/list", handlers.BuildHandleList(db))
	server.POST("/email/upload", handlers.BuildHandlerUpload(saveFileUseCase))
	server.GET("/email/:id", handlers.BuildHandleDetails(db))

	server.Run()
}
