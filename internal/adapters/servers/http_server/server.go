package http_server

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"messaggio/internal/queue"
	"messaggio/internal/utils"
	"net/http"
)

type HttpResponse struct {
	status int
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type MessageHttpServiceServer struct {
	amqpChannel *amqp.Channel
}

func StartServer() {


	appConf := utils.GetAppConfig()

	conn, err := amqp.Dial(appConf.AmpqHost)
	queue.HandleError(err, "Can't connect to AMQP")

	defer conn.Close()

	amqpChannel, err := conn.Channel()
	queue.HandleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	fmt.Println("Starting Message Service Server")

	mux := http.NewServeMux()

	serviceServer := MessageHttpServiceServer{amqpChannel: amqpChannel}

	mux.HandleFunc("/add_message", serviceServer.AddMessageHandler)

	err = http.ListenAndServe(appConf.Host, mux)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}