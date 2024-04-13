package main

import (
	"log"

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

	channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello, World!"),
		},
	)
}
