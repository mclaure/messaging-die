package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mclaure/messaging-die/appconfig"
	"github.com/mclaure/messaging-die/util"
	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
)

func failOnClientError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendRequestToRPCQueue(apiPattern string) (res []byte, err error) {
	conn, ch := util.GetChannel(
		appconfig.RabbitMQURL,
		appconfig.RabbitMQUsername,
		appconfig.RabbitMQPassword,
	)
	defer conn.Close()
	defer ch.Close()

	q := util.GetExclusiveQueue(ch)

	msgs := util.ConsumeFromExclusiveQueue(q.Name, ch)

	corrID := uuid.NewRandom().String()

	log.Printf("[C] CorrelationId (%s)", corrID)

	err = util.PublishToQueue(
		appconfig.RabbitMQRpcQueue,
		ch,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          []byte(apiPattern),
		},
	)

	// The response is handled here
	for d := range msgs {
		if corrID == d.CorrelationId {
			res = d.Body
			break
		}
	}

	return
}

func main() {
	startServer()
}

func startServer() {
	fmt.Println("Starting Public Service")
	router := mux.NewRouter()

	handleRequests(router)
	setupServer(router)
}

func handleRequests(router *mux.Router) {
	// GET Requests
	publicAPIPatterns := appconfig.GetPublicAPIPatterns()
	for _, publicAPIPattern := range publicAPIPatterns {
		handleGetRequest(router, publicAPIPattern)
	}

	// POST Requests
	// TODO
}

func setupServer(router *mux.Router) {
	srv := &http.Server{
		Handler:      router,
		Addr:         appconfig.PublicAPIFullAddress,
		WriteTimeout: appconfig.PublicAPIWriteTimeout * time.Second,
		ReadTimeout:  appconfig.PublicAPIReadTimeout * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handleGetRequest(router *mux.Router, apiPattern string) {
	router.HandleFunc(
		apiPattern,
		func(writer http.ResponseWriter, request *http.Request) {
			log.Printf("[C] Requesting APIPattern(%s)", apiPattern)

			res, err := sendRequestToRPCQueue(apiPattern)

			failOnClientError(err, "Failed to handle RPC Request")

			writer.Header().Set("Content-Type", "application/json")
			writer.Write(res)
		})
}
