package http_server

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	appqueue "messaggio/internal/queue"
	"net/http"
)

func (s *MessageHttpServiceServer) AddMessageHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	message, err := prepareRequestData(r)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		result := &HttpResponse{http.StatusBadRequest, "", err.Error()}
		showResponse(result, w)
		return
	}

	queue, err := s.amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	appqueue.HandleError(err, "Could not declare `add` queue")


	body, err := json.Marshal(message)

	if err != nil {
		result := &HttpResponse{http.StatusBadRequest, "", err.Error()}
		showResponse(result, w)
		return
	}


	err = s.amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("Message send: %s, %s", message.Phone, message.Text)
	result := &HttpResponse{http.StatusOK, "message send to quote", ""}

	showResponse(result, w)
}
