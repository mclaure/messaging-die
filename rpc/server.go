package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mclaure/messaging-die/appconfig"
	DieHttp "github.com/mclaure/messaging-die/http"
	"github.com/mclaure/messaging-die/util"
	"github.com/streadway/amqp"
)

var apiPatternMap = appconfig.GetAPIPatternMap()

func failOnServerError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getRealAPIURL(body []byte) string {
	var realAPIURL string

	requestInfo := DieHttp.GetRequestInfo(body)
	if len(requestInfo.URLParameters) > 0 {
		realAPIURL = fmt.Sprintf(
			"%s%s%s?%s",
			appconfig.DieAPIUrl,
			appconfig.DieAPIAddress,
			apiPatternMap[requestInfo.APIPattern],
			DieHttp.GetURLParameters(requestInfo.URLParameters),
		)
	} else {
		realAPIURL = fmt.Sprintf(
			"%s%s%s",
			appconfig.DieAPIUrl,
			appconfig.DieAPIAddress,
			apiPatternMap[requestInfo.APIPattern],
		)
	}

	return realAPIURL
}

func sendGetRequest(body []byte) []byte {
	realAPIURL := getRealAPIURL(body)
	log.Printf("[S] Send GetRequest: [%s]", realAPIURL)

	resp, err := http.Get(realAPIURL)
	defer resp.Body.Close()

	failOnServerError(err, fmt.Sprintf("Failed to get the response from: %s", realAPIURL))

	data, err := ioutil.ReadAll(resp.Body)
	failOnServerError(err, fmt.Sprintf("Failed to read the response from: %s", realAPIURL))

	return data
}

func main() {
	conn, ch := util.GetChannel(
		appconfig.RabbitMQURL,
		appconfig.RabbitMQUsername,
		appconfig.RabbitMQPassword,
	)
	defer conn.Close()
	defer ch.Close()

	q := util.GetRPCQueue(appconfig.RabbitMQRpcQueue, ch)

	msgs := util.ConsumeFromRPCQueue(q.Name, ch)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			response := sendGetRequest(d.Body)

			err := util.PublishToQueue(
				d.ReplyTo,
				ch,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          response,
				},
			)
			failOnServerError(err, fmt.Sprintf("Failed to publish to queue: %s", d.ReplyTo))

			d.Ack(false)
		}
	}()

	log.Printf("[S] Waiting for RPC requests")
	<-forever
}
