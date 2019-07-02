package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

var rabbitConn *amqp.Connection
var rabbitCh *amqp.Channel
var rabbitErrCh chan *amqp.Error

const rabbitConnStr = "amqp://rabbitmq:rabbitmq@localhost:5672/" // can be arg/option
const rabbitExchange = ""                                        // can be arg/option
const rabbitQueueName = "hello"                                  // can be arg/option

func main() {
	defer closePreviousRabbit()

	done := make(chan bool)
	rabbitErrCh = make(chan *amqp.Error)

	go func() {
		for {
			amqpErr := <-rabbitErrCh
			if nil != amqpErr {
				log.Println("Error chan:", amqpErr)
			} else {
				closePreviousRabbit()
			}

			rabbitErrCh = make(chan *amqp.Error)

			log.Println("try reconnect")

			var err error

			rabbitConn, err = amqp.Dial(rabbitConnStr)
			if nil != err {
				log.Printf("Failed to connect to RabbitMQ: %v", err)
				continue
			}
			rabbitConn.NotifyClose(rabbitErrCh)

			rabbitCh, err = rabbitConn.Channel()
			if nil != err {
				log.Printf("Failed to open a channel: %v", err)
				continue
			}

			log.Println("connection was opened")
		}
	}()

	body := "Hello World!"
	counter := 0

	go func() {
		for {
			if nil == rabbitConn || nil == rabbitCh {
				log.Println("connection is not ready")
				rabbitErrCh <- nil
				time.Sleep(5 * time.Second)
			} else {
				q, err := rabbitCh.QueueDeclare(
					rabbitQueueName, // name
					false,           // durable
					false,           // delete when unused
					false,           // exclusive
					false,           // no-wait
					nil,             // arguments
				)
				if nil != err {
					log.Printf("Queue Declare error: %v", err)
					rabbitErrCh <- nil
					time.Sleep(5 * time.Second)
					continue
				}
				err = rabbitCh.Publish(
					rabbitExchange, // exchange
					q.Name,         // routing key
					false,          // mandatory
					false,          // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(body),
					})

				if nil == err {
					log.Printf(" [x] Sent %s [%d]", body, counter)
					counter++
					time.Sleep(3 * time.Second)
				} else {
					log.Printf("Failed to publish a message: %v", err)
					rabbitErrCh <- nil
					time.Sleep(5 * time.Second)
					continue
				}

			}
		}
	}()

	rabbitErrCh <- nil

	<-done

}

func closePreviousRabbit() {
	if nil != rabbitConn {
		rabbitConn.Close()
	}
	if nil != rabbitCh {
		rabbitCh.Close()
	}

	rabbitConn = nil
	rabbitCh = nil
}
