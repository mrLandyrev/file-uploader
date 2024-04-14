package handlers

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

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
		"null":       "id",
	}

	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		perPage, _ := strconv.Atoi(c.Query("perPage"))
		sortBy := sortByKeyMap[c.Query("sortBy")]
		sortDirection := c.Query("sortDirection")
		target := c.QueryArray("regExps")

		cr := make([]regexp.Regexp, 0, len(target))

		for _, t := range target {
			r, err := regexp.Compile(t)
			if err != nil {
				c.Status(400)
				c.Done()
			}
			cr = append(cr, *r)
		}

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
				content,
				(select COUNT(id) from email)
			from email order by "%s" %s
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
				content     string
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
				&content,
				&count,
			)

			dangerValues := make([]string, 0)

			for _, r := range cr {
				dangerValues = append(dangerValues, r.FindAllString(string(content), -1)...)
			}

			res = append(res, gin.H{
				"id":           id,
				"fileName":     filename,
				"sendDate":     send_date,
				"uploadDate":   upload_date,
				"to":           strings.Split(to, ","),
				"cc":           strings.Split(cc, ","),
				"from":         strings.Split(from, ","),
				"subject":      subject,
				"isDanger":     len(dangerValues) > 0,
				"dangerValues": dangerValues,
			})
		}

		resCount := count / perPage
		if perPage*resCount < count {
			resCount++
		}

		c.JSON(200, gin.H{
			"items":      res,
			"totalItems": count,
			"pageCount":  resCount,
		})
	}
}
