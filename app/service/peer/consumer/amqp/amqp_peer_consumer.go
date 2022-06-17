package amqp

import "github.com/streadway/amqp"



type AmqpPeerConsumer struct {
	connection *amqp.Connection
}