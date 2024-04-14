package usecases

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"time"
)

type ParsedResult struct {
	SendDate time.Time
	To       []string
	Cc       []string
	Subject  string
	Content  io.ReadWriter
	From     []string
}

func BuildFileName(part *multipart.Part, radix string, index int) (filename string) {

	filename = part.FileName()
	if len(filename) > 0 {
		return
	}

	mediaType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err == nil {
		mime_type, e := mime.ExtensionsByType(mediaType)
		if e == nil {
			ext := ".bin"
			if len(mime_type) > 0 {
				ext = mime_type[0]
			}
			return fmt.Sprintf("%s-%d%s", radix, index, ext)
		}
	}

	return

}

func WritePart(part *multipart.Part, filename string, resMap map[string]io.ReadWriter) {

	part_data, err := io.ReadAll(part)
	if err != nil {
		log.Println("Error reading MIME part data -", err)
		return
	}

	content_transfer_encoding := strings.ToUpper(part.Header.Get("Content-Transfer-Encoding"))

	switch {

	case strings.Compare(content_transfer_encoding, "BASE64") == 0:
		decoded_content, err := base64.StdEncoding.DecodeString(string(part_data))
		if err != nil {
			log.Println("Error decoding base64 -", err)
		} else {
			_, ok := resMap[filename]
			if !ok {
				resMap[filename] = new(bytes.Buffer)
			}
			resMap[filename].Write(decoded_content)
		}

	case strings.Compare(content_transfer_encoding, "QUOTED-PRINTABLE") == 0:
		decoded_content, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(part_data)))
		if err != nil {
			log.Println("Error decoding quoted-printable -", err)
		} else {
			_, ok := resMap[filename]
			if !ok {
				resMap[filename] = new(bytes.Buffer)
			}
			resMap[filename].Write(decoded_content)
		}

	default:
		_, ok := resMap[filename]
		if !ok {
			resMap[filename] = new(bytes.Buffer)
		}
		resMap[filename].Write(part_data)

	}

}

func ParsePart(mime_data io.Reader, boundary string, index int, resMap map[string]io.ReadWriter) {

	reader := multipart.NewReader(mime_data, boundary)
	if reader == nil {
		return
	}

	for {

		new_part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		mediaType, params, err := mime.ParseMediaType(new_part.Header.Get("Content-Type"))
		if err == nil && strings.HasPrefix(mediaType, "multipart/") {
			ParsePart(new_part, params["boundary"], index+1, resMap)
		} else {
			filename := BuildFileName(new_part, boundary, 1)
			WritePart(new_part, filename, resMap)
		}

	}

}

func ParseFiles(f io.Reader) ParsedResult {

	m, err := mail.ReadMessage(f)
	if err != nil {
		log.Fatalln("Parse mail KO -", err)
	}

	res := &ParsedResult{}

	dec := new(mime.WordDecoder)
	from, _ := dec.DecodeHeader(m.Header.Get("From"))
	to, _ := dec.DecodeHeader(m.Header.Get("To"))
	cc, _ := dec.DecodeHeader(m.Header.Get("Cc"))
	subject, _ := dec.DecodeHeader(m.Header.Get("Subject"))
	sendDate, _ := dec.DecodeHeader(m.Header.Get("Date"))

	res.To = strings.Split(to, ", ")
	res.Cc = strings.Split(cc, ", ")
	res.Subject = subject
	res.From = strings.Split(from, ", ")
	res.SendDate, _ = time.Parse("2 Jan 2006 15:04:05 -0700", sendDate)
	res.Content = new(bytes.Buffer)

	mediaType, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		content_transfer_encoding, _ := dec.DecodeHeader(m.Header.Get("Content-Transfer-Encoding"))
		data, _ := io.ReadAll(m.Body)
		switch {

		case strings.Compare(content_transfer_encoding, "base64") == 0:
			decoded_content, err := base64.StdEncoding.DecodeString(string(data))
			if err != nil {
				log.Println("Error decoding base64 -", err)
			} else {
				res.Content.Write(decoded_content)
			}

		case strings.Compare(content_transfer_encoding, "quoted-printable") == 0:
			decoded_content, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(data)))
			if err != nil {
				log.Println("Error decoding quoted-printable -", err)
			} else {
				res.Content.Write(decoded_content)
			}

		default:
			res.Content.Write(data)

		}
		return *res
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		content_transfer_encoding, _ := dec.DecodeHeader(m.Header.Get("Content-Transfer-Encoding"))
		data, _ := io.ReadAll(m.Body)
		switch {

		case strings.Compare(content_transfer_encoding, "base64") == 0:
			decoded_content, err := base64.StdEncoding.DecodeString(string(data))
			if err != nil {
				log.Println("Error decoding base64 -", err)
			} else {
				res.Content.Write(decoded_content)
			}

		case strings.Compare(content_transfer_encoding, "quoted-printable") == 0:
			decoded_content, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(data)))
			if err != nil {
				log.Println("Error decoding quoted-printable -", err)
			} else {
				res.Content.Write(decoded_content)
			}

		default:
			res.Content.Write(data)

		}
		return *res
	}

	resMap := make(map[string]io.ReadWriter, 0)
	ParsePart(m.Body, params["boundary"], 1, resMap)
	resBuf := new(bytes.Buffer)
	for key, value := range resMap {
		if strings.HasSuffix(key, ".bin") {
			continue
		}

		io.Copy(resBuf, value)
	}

	res.Content = resBuf
	return *res
}
