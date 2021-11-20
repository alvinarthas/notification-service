package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	// Will Send email body or sms body to the broker
	// Create Connection to the rabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the Exchange
	err = ch.ExchangeDeclare(
		"notification", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Declare Queue for SMS
	qs, err := ch.QueueDeclare(
		"sms-queue", // name
		false,       // durable
		false,       // delete when unused
		true,        // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue for sms")

	// Bind the Queue with the Exchange
	err = ch.QueueBind(
		qs.Name,        // queue name
		"sms",          // routing key
		"notification", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a sms queue")

	smsConsume, err := ch.Consume(
		qs.Name, // queue
		"",      // consumer
		true,    // auto ack
		false,   // exclusive
		false,   // no local
		false,   // no wait
		nil,     // args
	)
	failOnError(err, "Failed to register a sms consumer")

	// Declare Queue for Mail
	qm, err := ch.QueueDeclare(
		"mail-queue", // name
		false,        // durable
		false,        // delete when unused
		true,         // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue for mail")

	// Bind the Queue with the Exchange
	err = ch.QueueBind(
		qm.Name,        // queue name
		"email",        // routing key
		"notification", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a mail queue")

	// Consume from Queue
	mailConsume, err := ch.Consume(
		qm.Name, // queue
		"",      // consumer
		true,    // auto ack
		false,   // exclusive
		false,   // no local
		false,   // no wait
		nil,     // args
	)
	failOnError(err, "Failed to register a mail consumer")

	// Send the information
	forever := make(chan bool)

	go func() {

		for sm := range smsConsume {
			log.Printf(" [x] - SMS Consume - %s", sm.Body)
		}

	}()

	go func() {

		for mm := range mailConsume {
			log.Printf(" [x] - Mail Consume - %s", mm.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
