package main

import (
	"log"
	"time"

	"github.com/goforbroke1006/fault-tolerance-samples/reconnect-on-fail-sync-lazy/pkg/reconnect"
)

const rabbitConnStr = "amqp://rabbitmq:rabbitmq@localhost:5672/" // can be arg/option
const rabbitExchange = ""                                        // can be arg/option
const rabbitQueueName = "hello"                                  // can be arg/option

func main() {
	rabbit := reconnect.NewRabbitMQReconnector(rabbitConnStr)
	defer rabbit.Close()

	done := make(chan bool)

	body := "Hello World!"
	counter := 0

	go func() {
		for {
			err := rabbit.Publish(rabbitExchange, rabbitQueueName, []byte(body))

			if nil == err {
				log.Printf(" [x] Sent %s [%d]", body, counter)
				counter++
				time.Sleep(3 * time.Second)
			} else {
				log.Printf("Failed to publish a message: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
		}
	}()

	<-done

}
