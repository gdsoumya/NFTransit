package queue

import "github.com/streadway/amqp"

type QConfig struct {
	User      string
	Password  string
	QueueName string
	Host      string
	Port      uint64
}

type IQueue interface {
	Close() error
	Channel() (*amqp.Channel, error)
	DeclareExchangeQueue() (amqp.Queue, error)
	Qos(prefectCount int) error
	Consume() (<-chan amqp.Delivery, error)
	SetupPublisherConfirms() (chan amqp.Confirmation, error)
	Publish(body string, delay int64) error
	ConfirmPublish() error
	NackDelivery(d *amqp.Delivery, multiple bool, requeue bool) error
}
