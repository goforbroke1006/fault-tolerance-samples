package reconnect

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
)

func NewRabbitMQReconnector(rabbitConnStr string) *RabbitMQReconnector {
	return &RabbitMQReconnector{
		rabbitConnStr: rabbitConnStr,
	}
}

type RabbitMQReconnector struct {
	rabbitConnStr string

	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel

	// rabbitErrCh chan *amqp.Error
}

func (rr *RabbitMQReconnector) IsReady() bool {
	return !(nil == rr.rabbitConn || rr.rabbitConn.IsClosed() || nil == rr.rabbitCh)
}

func (rr *RabbitMQReconnector) Publish(exchangeName, queueName string, body []byte) error {
	if !rr.IsReady() {
		log.Println("connection is not ready")
		rr.checkAndFix()

		if !rr.IsReady() {
			return errors.New("connection is not ready")
		} else {
			return rr.publish(exchangeName, queueName, body)
		}
	}

	return rr.publish(exchangeName, queueName, body)
}

func (rr RabbitMQReconnector) Close() error {
	rr.closePreviousRabbit()
	return nil
}

func (rr *RabbitMQReconnector) closePreviousRabbit() {
	if nil != rr.rabbitConn {
		rr.rabbitConn.Close()
	}
	if nil != rr.rabbitCh {
		rr.rabbitCh.Close()
	}

	rr.rabbitConn = nil
	rr.rabbitCh = nil
}

func (rr *RabbitMQReconnector) checkAndFix() {
	log.Println("try reconnect")

	var err error

	conn, err := amqp.Dial(rr.rabbitConnStr)
	if nil != err {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	// conn.NotifyClose(rr.rabbitErrCh)
	rr.rabbitConn = conn

	ch, err := conn.Channel()
	if nil != err {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	rr.rabbitCh = ch

	log.Println("connection was opened")
}

func (rr *RabbitMQReconnector) publish(exchangeName, queueName string, body []byte) error {
	q, err := rr.rabbitCh.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if nil != err {
		return err
	}
	err = rr.rabbitCh.Publish(
		exchangeName, // exchange
		q.Name,       // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if nil != err {
		return err
	}

	return nil
}
