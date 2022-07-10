package consumer

import (
	"broker/core/config"
	"broker/core/dto"
	"context"
	"sync"

	"github.com/streadway/amqp"
)

type Consumer struct {
	connection            *amqp.Connection
	sendDelivery          chan dto.PeerEnvelope
	senderDispatchers     map[string]*SenderDispatcher
	senderDispatchersLock sync.RWMutex
}

func NewConsumer(connection *amqp.Connection) *Consumer {
	return &Consumer{
		connection:        connection,
		senderDispatchers: make(map[string]*SenderDispatcher, 1000),
		sendDelivery:      make(chan dto.PeerEnvelope),
	}
}

func (c *Consumer) Init() error {
	channel, err := c.connection.Channel()
	if err != nil {
		return err
	}

	if err := channel.ExchangeDeclare(
		getDeadLetterExchangeName(),
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := channel.QueueDeclare(
		getDeadLetterQueueName(),
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := channel.QueueBind(
		getDeadLetterQueueName(),
		"",
		getDeadLetterExchangeName(),
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateConnectionBase, error) {
	c.senderDispatchersLock.Lock()
	defer c.senderDispatchersLock.Unlock()

	dispatcher, exist := c.senderDispatchers[payload.WorkspaceID]
	if !exist {
		dispatcher = NewSenderDispatcher(c.connection, c.sendDelivery)
		c.senderDispatchers[payload.WorkspaceID] = dispatcher
	}
	queue, err := dispatcher.AddQueue(senderID, payload.WorkspaceID)
	if err != nil {
		return nil, err
	}

	config := config.GetConfig().AMQP
	return &dto.CreateConnectionBase{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.ExternalUser,
		Password: config.ExternalPassword,
		Vhost:    config.Vhost,
		Meta: dto.Meta{
			QueueName:    queue.Name(),
			ExchangeName: queue.ExchangeName(),
		},
	}, nil
}

func (c *Consumer) Consume(handler func(payload dto.PeerEnvelope)) {
	go func() {
		for payload := range c.sendDelivery {
			go handler(payload)
		}
	}()
}
