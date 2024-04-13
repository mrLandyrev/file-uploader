package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sg3des/eml"
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

	server.GET("/email/list", BuildHandleList(db))

	server.POST("/email/upload", func(c *gin.Context) {
		var file *multipart.FileHeader
		var err error
		if file, err = c.FormFile("email"); err != nil {
			c.Status(500)
			c.Done()
		}

		f, _ := file.Open()
		fileBytes, _ := io.ReadAll(f)
		email, _ := eml.Parse(fileBytes)

		To := ""
		for i, address := range email.To {
			if i > 0 {
				To += ","
			}
			To += address.Email()
		}

		Cc := ""
		for i, address := range email.Cc {
			if i > 0 {
				Cc += ","
			}
			Cc += address.Email()
		}

		From := ""
		for i, address := range email.From {
			if i > 0 {
				From += ","
			}
			From += address.Email()
		}

		content, _ := b64.StdEncoding.DecodeString(email.Text)
		subject, _ := b64.StdEncoding.DecodeString(email.Subject)

		rows, err := db.Query(`
			INSERT INTO email (
				id,
				filename,
				send_date,
				upload_date,
				"to",
				cc,
				subject,
				content,
				"from"
			) VALUES (
				gen_random_uuid (), 
				$1,
				$2,
				NOW (),
				$3,
				$4,
				$5,
				$6,
				$7
			)`,
			file.Filename,
			email.Date,
			To,
			Cc,
			subject,
			content,
			From,
		)

		if err != nil {
			log.Fatal("Database error", err)
			c.Status(500)
			c.Done()
		}
		defer rows.Close()

		c.Status(200)
		c.Done()
	})

	server.Run()
}

func BuildHandleList(db *sql.DB) func(c *gin.Context) {
	sortByKeyMap := map[string]string{
		"id":         "id",
		"fileName":   "filename",
		"sendDate":   "send_date",
		"uploadDate": "upload_date",
		"to":         "to",
		"cc":         "cc",
		"from":       "from",
		"subject":    "subject",
	}

	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		perPage, _ := strconv.Atoi(c.Query("perPage"))
		sortBy := sortByKeyMap[c.Query("sortBy")]
		sortDirection := c.Query("sortDirection")

		q := fmt.Sprintf(`
			select
				id,
				filename,
				send_date,
				upload_date,
				"to",
				cc,
				"from",
				subject,
				(select COUNT(id) from email)
			from email order by %s %s
				limit $1
				offset $2
		`, sortBy, sortDirection)

		rows, _ := db.Query(
			q,
			perPage,
			page*perPage,
		)
		defer rows.Close()

		res := make([]gin.H, 0)
		count := 0

		for rows.Next() {
			var (
				id          string
				filename    string
				send_date   time.Time
				upload_date time.Time
				to          string
				cc          string
				from        string
				subject     string
			)

			rows.Scan(
				&id,
				&filename,
				&send_date,
				&upload_date,
				&to,
				&cc,
				&from,
				&subject,
				&count,
			)

			res = append(res, gin.H{
				"id":         id,
				"fileName":   filename,
				"sendDate":   send_date,
				"uploadDate": upload_date,
				"to":         strings.Split(to, ","),
				"cc":         strings.Split(cc, ","),
				"from":       strings.Split(from, ","),
				"subject":    subject,
			})
		}

		resCount := count / perPage
		if perPage*resCount < count {
			resCount++
		}

		c.JSON(200, gin.H{
			"items":     res,
			"pageCount": resCount,
		})
	}
}
