package main

import (
	"log"
	"os"

	"github.com/mrLandyrev/file-uploader/internal/loader"
)

func main() {
	appConfig := &loader.Config{
		Path: os.Args[1],
	}
	app, err := loader.NewApp(*appConfig)

	if err != nil {
		log.Fatalln("Error at creating app. Stop.", err)
	}

	app.Run()
}
