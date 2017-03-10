package loggerclient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

var (
	uri               = ""
	queue             = "logs"
	conn              *amqp.Connection
	ch                *amqp.Channel
	q                 *amqp.Queue
	dialRetryCount    = 0
	maxConnRetryCount int
)

// LoggerClient is the strucutre of the logger NewClient
// You can use this logger client to send logs ot the fastbase logger service
type LoggerClient struct {
	LogSet  string
	LogType string
}

func init() {
	setupEnvVars()
}

// NewClient creates a new client and predefined Log Set and logTypego bu
// logSet Is the default log set (elastic index)
// logType Is the type of the log (elastic type)
func NewClient(logSet string, logType string) *LoggerClient {
	conn = dial()
	ch = createChannel(conn)
	q = createQueue(queue, ch)
	return &LoggerClient{LogSet: logSet, LogType: logType}
}

// LogWithType allows you to log and override default logType paramter of the logger client
func (c *LoggerClient) LogWithType(logType string, log map[string]interface{}) error {
	body := c.createMessageBody(logType, log)

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

// Log allows you to log a message with default logType
func (c *LoggerClient) Log(log map[string]interface{}) error {
	return c.LogWithType(c.LogType, log)
}

func (c *LoggerClient) createMessageBody(logType string, log map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"log_set":  c.LogSet,
		"log_type": logType,
		"log":      log,
	}
}

func dial() *amqp.Connection {
	conn, err := amqp.Dial(uri)
	if err != nil && dialRetryCount < maxConnRetryCount {
		log.Printf("Failed to connect Rabbit MQ at %s, will retry.....%d", uri, dialRetryCount+1)
		dialRetryCount++
		time.Sleep(time.Second * 1)
		return dial()
	}
	failOnError(err, "Failed to setup Rabbit MQ connection")
	dialRetryCount = 0
	log.Printf("Connected to Rabbit MQ with uri: %s", uri)
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

func getEnvOrFail(key string) string {
	env := os.Getenv(key)

	if env == "" {
		panic(key)
	}
	return env
}

func setupEnvVars() {
	uri = getEnvOrFail("MQ_CONN_STRING")
	count, err := strconv.Atoi(getEnvOrFail("CONN_RETRY_COUNT"))
	maxConnRetryCount = count
	failOnError(err, "CONN_RETRY_COUNT must be an integer type")
}
