package usecases

import (
	"database/sql"
	"io"
	"log"
)

func BuildSaveFileUseCase(db *sql.DB) func(parsedResult ParsedResult, filename string) error {
	return func(parsedResult ParsedResult, filename string) error {
		To := ""
		for i, address := range parsedResult.To {
			if i > 0 {
				To += ","
			}
			To += address
		}

		Cc := ""
		for i, address := range parsedResult.Cc {
			if i > 0 {
				Cc += ","
			}
			Cc += address
		}

		From := ""
		for i, address := range parsedResult.From {
			if i > 0 {
				From += ","
			}
			From += address
		}

		content, _ := io.ReadAll(parsedResult.Content)

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
			filename,
			parsedResult.SendDate,
			To,
			Cc,
			parsedResult.Subject,
			content,
			From,
		)

		if err != nil {
			log.Fatal("Database error", err)
			return err
		}
		defer rows.Close()

		return nil
	}
}
