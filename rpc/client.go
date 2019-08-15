package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mclaure/messaging-die/appconfig"
	DieHttp "github.com/mclaure/messaging-die/http"
	"github.com/mclaure/messaging-die/util"
	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
)

/***********************************************************************
 *                          Start ClientServer                         *
 ***********************************************************************/
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
	publicGetAPIPatterns := appconfig.GetPublicGetAPIPatterns()
	for _, publicGetAPIPattern := range publicGetAPIPatterns {
		handleGetRequest(router, publicGetAPIPattern)
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

func failOnClientError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

/***********************************************************************
 *                         Handle Get Requests                         *
 ***********************************************************************/
func handleGetRequest(router *mux.Router, apiPattern string) {
	router.HandleFunc(apiPattern, delegateGetRequest)
}

func delegateGetRequest(writer http.ResponseWriter, request *http.Request) {
	res, err := sendRequestToRPCQueue(request)

	failOnClientError(err, "Failed to handle RPC Request")

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(res)
}

/***********************************************************************
 *                         Handle Post Requests                        *
 ***********************************************************************/
func handlePostRequest(router *mux.Router, apiPattern string) {
	router.HandleFunc(apiPattern, delegatePostRequest)
}

func delegatePostRequest(writer http.ResponseWriter, request *http.Request) {
	// TODO
}

/***********************************************************************
 *                      Send Request To RPC Queue                      *
 ***********************************************************************/
func sendRequestToRPCQueue(request *http.Request) (res []byte, err error) {
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
			Body:          DieHttp.GetRequestInfoBytes(request),
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
