package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	uri   = "amqp://guest:guest@localhost:5672/"
	queue = "logs"
	conn  *amqp.Connection
	ch    *amqp.Channel
	q     *amqp.Queue
)

func main() {}

type LoggerClient struct {
	LogSet  string
	LogType string
}

func NewClient(logSet string, logType string) *LoggerClient {
	conn = dial()
	ch = createChannel(conn)
	q = createQueue(queue, ch)
	return &LoggerClient{LogSet: logSet, LogType: logType}
}

func (c *LoggerClient) LogWithType(logType string, log map[string]interface{}) error {
	body := map[string]interface{}{
		"log_set":  c.LogSet,
		"log_type": logType,
		"log":      log,
	}

	data, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})
	return err
}

func (c *LoggerClient) Log(log map[string]interface{}) error {
	return c.LogWithType(c.LogType, log)
}

func dial() *amqp.Connection {
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to setup Rabbit MQ connection")
	return conn
}

func createChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to create a channel")
	return ch
}

func createQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(name, true, false, false, false, nil)
	failOnError(err, "Faield to declare a queue")

	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set a prefetch policy")

	return &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
