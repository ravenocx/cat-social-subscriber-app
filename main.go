package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ : %+v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel : %+v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cat_matches", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare a queue : %+v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("Failed to register a consumer : %+v", err)
	}

	var forever chan struct{}
	
	go func(){
		for{
			if len(msgs) == 0 {
				log.Printf(" [*] Waiting for notification......")
				time.Sleep(5 * time.Second)
			}
		}
	}()
	go func() {
		for d := range msgs {
			log.Printf("Received a Notification matching cats: %s", d.Body)
		}
	}()
	
	<-forever
}
