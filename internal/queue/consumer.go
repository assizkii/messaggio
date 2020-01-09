package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"messaggio/internal/domain/interfaces"
	"os"
)

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

func Serve(host string) *amqp.Channel  {

	conn, err := amqp.Dial(host)
	HandleError(err, "Can't connect to AMQP")


	defer conn.Close()

	amqpChannel, err := conn.Channel()
	HandleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	return amqpChannel
}

func RunConsumer(amqpChannel *amqp.Channel, storage interfaces.MessageStorage) {

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	HandleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	HandleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	HandleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			message := &interfaces.Message{}

			err := json.Unmarshal(d.Body, message)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			newId, err := storage.Add(message)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("New message added to database with id %d", newId)
			}

		}
	}()

	// Stop for program termination
	<-stopChan

}