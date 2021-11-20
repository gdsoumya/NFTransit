package queue

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gdsoumya/nftransit/relay/pkg/env"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Queue struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	confirms chan amqp.Confirmation
	logger   *zap.Logger
	config   *QConfig
}

func NewQueueClient(logger *zap.Logger, conf *QConfig) (*Queue, error) {
	// create connection
	dsn := url.URL{
		User:   url.UserPassword(conf.User, conf.Password),
		Scheme: "amqp",
		Host:   fmt.Sprintf("%s:%v", conf.Host, conf.Port),
	}

	conn, err := amqp.Dial(dsn.String())
	if err != nil {
		logger.Debug("failed to connect to rabbitmq", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to rabbitmq, error: %w", err)
	}

	return &Queue{
		conn:   conn,
		logger: logger,
		config: conf,
	}, nil
}

func (q *Queue) Close() error {
	return q.conn.Close()
}

func (q *Queue) Channel() (*amqp.Channel, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		return nil, err
	}
	q.ch = ch
	return ch, nil
}

func (q *Queue) DeclareExchangeQueue() (amqp.Queue, error) {
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	err := q.ch.ExchangeDeclare("delayed", "x-delayed-message", true, false, false, false, args)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare delayed exchange, error:%w", err)
	}

	queue, err := q.ch.QueueDeclare(
		q.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue, error:%w", err)
	}

	return queue, q.ch.QueueBind(queue.Name, queue.Name, "delayed", false, nil)
}

func (q *Queue) Qos(prefectCount int) error {
	return q.ch.Qos(prefectCount, 0, false)
}

func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	return q.ch.Consume(q.config.QueueName, "", false, false, false, false, nil)
}

func (q *Queue) SetupPublisherConfirms() (chan amqp.Confirmation, error) {
	// Buffer of 1 for our single outstanding publishing
	q.confirms = q.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	if err := q.ch.Confirm(false); err != nil {
		return nil, err
	}
	return q.confirms, nil
}

func (q *Queue) Publish(body string, delay int64) error {
	headers := make(amqp.Table)
	if delay != 0 {
		headers["x-delay"] = delay
	}
	return q.ch.Publish(
		"delayed",
		q.config.QueueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "application/json",
			Body:         []byte(body),
			Headers:      headers,
		})
}

func (q *Queue) AckDelivery(d *amqp.Delivery, multiple bool) error {
	return d.Ack(multiple)
}

func (q *Queue) NackDelivery(d *amqp.Delivery, multiple bool, requeue bool) error {
	return d.Nack(multiple, requeue)
}

func (q *Queue) ConfirmPublish() error {
	if confirmed := <-q.confirms; !confirmed.Ack {
		return fmt.Errorf("publisher confirm failed")
	}
	return nil
}

func InitQueue(logger *zap.Logger, envConfig *env.EnvData, publisher bool) (*Queue, error) {
	queueClient, err := NewQueueClient(logger, &QConfig{
		User:      envConfig.QueueUser,
		Password:  envConfig.QueuePassword,
		QueueName: envConfig.QueueName,
		Host:      envConfig.QueueHost,
		Port:      envConfig.QueuePort,
	})
	if err != nil {
		return nil, err
	}
	if _, err = queueClient.Channel(); err != nil {
		return nil, err
	}
	if _, err = queueClient.DeclareExchangeQueue(); err != nil {
		return nil, err
	}
	if publisher {
		if _, err = queueClient.SetupPublisherConfirms(); err != nil {
			return nil, err
		}
	} else {
		if err = queueClient.Qos(1); err != nil {
			return nil, err
		}
	}
	return queueClient, nil
}
