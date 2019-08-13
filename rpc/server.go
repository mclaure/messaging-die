package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mclaure/messaging-die/appconfig"
	"github.com/mclaure/messaging-die/util"
	"github.com/streadway/amqp"
)

func failOnServerError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getAPIURL(pattern string) string {
	// We can use a map here, but golang doesn't allow us to hava a const map
	if pattern == appconfig.PublicAPICountriesPattern {
		pattern = appconfig.DieAPICountriesPattern
	}

	if pattern == appconfig.PublicAPITeamsPattern {
		pattern = appconfig.DieAPITeamsPattern
	}

	return fmt.Sprintf("%s%s%s", appconfig.DieAPIUrl, appconfig.DieAPIAddress, pattern)
}

func sendGetRequest(apiPattern string) []byte {
	apiURL := getAPIURL(apiPattern)
	log.Printf("[S] Send GetRequest: [%s]", apiURL)

	resp, err := http.Get(apiURL)
	defer resp.Body.Close()

	failOnServerError(err, fmt.Sprintf("Failed to get the response from: %s", apiURL))

	data, err := ioutil.ReadAll(resp.Body)
	failOnServerError(err, fmt.Sprintf("Failed to read the response from: %s", apiURL))

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
			response := sendGetRequest(string(d.Body))

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
