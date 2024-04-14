package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func BuildHandleDetails(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		paramId := c.Param("id")

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
				content
			from email where id = '%s'
		`, paramId)

		row := db.QueryRow(q)

		if row.Err() != nil {
			c.Status(500)
			c.Done()
		}

		var (
			id          string
			filename    string
			send_date   time.Time
			upload_date time.Time
			to          string
			cc          string
			from        string
			subject     string
			content     string
		)

		row.Scan(
			&id,
			&filename,
			&send_date,
			&upload_date,
			&to,
			&cc,
			&from,
			&subject,
			&content,
		)

		c.JSON(200, gin.H{
			"id":         id,
			"fileName":   filename,
			"sendDate":   send_date,
			"uploadDate": upload_date,
			"to":         strings.Split(to, ","),
			"cc":         strings.Split(cc, ","),
			"from":       strings.Split(from, ","),
			"subject":    subject,
			"content":    content,
		})
	}
}
