package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type Mail struct {
	Receiver []string
	Sender   string
	Subject  string
	MailBody string
}

type SMS struct {
	Receiver string
	Sender   string
	Message  string
}

const (
	EMAIL_NOTIFICATION = "email"
	SMS_NOTIFICATION   = "sms"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	var (
		err              error
		notificationType = os.Args[1]
		body             []byte
	)

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

	// Check the input type
	switch notificationType {
	case EMAIL_NOTIFICATION:
		body, err = GenerateMailBody()
		failOnError(err, "Failed to generate body email")
	case SMS_NOTIFICATION:
		body, err = GenerateSMSBody()
		failOnError(err, "Failed to generate body sms")
	default:
		failOnError(err, "wrong input type!")
		return
	}

	// Publish to the Exchange
	err = ch.Publish(
		"notification",   // exchange
		notificationType, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("Successfully Sent %s", notificationType)
}

func GenerateMailBody() ([]byte, error) {

	// Create the Body
	receiver := []string{"alvin@mailinator.com", "joni@mailinator.com"}
	sender := "rebitemki@mailinator.com"
	subject := "Test Send Email"
	body := "Hai Welcome to Our Community"

	mail := Mail{
		Receiver: receiver,
		Sender:   sender,
		Subject:  subject,
		MailBody: body,
	}
	request, err := json.Marshal(mail)

	return request, err
}

func GenerateSMSBody() ([]byte, error) {

	// Create the Body
	receiver := "0811230329"
	sender := "0939842389"
	message := "Hai Welcome to Our Community"

	mail := SMS{
		Receiver: receiver,
		Sender:   sender,
		Message:  message,
	}

	request, err := json.Marshal(mail)

	return request, err
}
