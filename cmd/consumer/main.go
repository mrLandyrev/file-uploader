package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ArkaGPL/parsemail"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@192.168.0.105:5672/")
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"raw", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	rawMessages, _ := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for message := range rawMessages {
		msg, _ := parsemail.Parse(bytes.NewReader(message.Body))
		fmt.Println(len(msg.Attachments))
		// str, err := io.ReadAll()
		// fmt.Println(msg.TextBody)
		// fmt.Println(len(msg.Attachments))

		// for i, attachment := range msg.Attachments {
		// 	fmt.Println(string(i))
		// 	at, _ := io.ReadAll(attachment.Data)
		// 	fmt.Println(string(at))
		// }

		// fmt.Println(string(str))

		// fmt.Println(string(message.Body))

		// b, _ := b64.RawStdEncoding.DecodeString(string(message.Body))

		// fmt.Println(string(b))
	}
}
