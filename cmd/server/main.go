package main

import (
	"log"

	"github.com/mrLandyrev/file-uploader/internal/server"
)

func main() {
	appConfig := &server.Config{}
	app, err := server.NewApp(*appConfig)

	if err != nil {
		log.Fatalln("Error at creating app. Stop.", err)
	}

	app.Run()
}
