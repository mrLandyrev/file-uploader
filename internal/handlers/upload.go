package handlers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/mrLandyrev/file-uploader/internal/usecases"
)

func BuildHandlerUpload(saveFileUseCase func(parsedResult usecases.ParsedResult, filename string) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		var file *multipart.FileHeader
		var err error
		if file, err = c.FormFile("email"); err != nil {
			c.Status(500)
			c.Done()
		}

		f, _ := file.Open()
		defer f.Close()

		err = saveFileUseCase(usecases.ParseVor(f), file.Filename)

		if err != nil {
			c.Status(500)
			c.Done()
		}

		c.Status(200)
		c.Done()
	}
}
