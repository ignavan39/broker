package amqp

import (
	"broker/app/config"
	"broker/app/dto"
	"context"
	"sync"

	"github.com/streadway/amqp"
)

type AmqpWorkspaceConsumer struct {
	connection           *amqp.Connection
	sendDelivery         chan dto.PeerEnvelope
	senderDispatchers    map[string]*SenderDispatcher
	senderDispatchersMux sync.RWMutex
}

func NewAmqpWorkspaceConsumer(connection *amqp.Connection) *AmqpWorkspaceConsumer {
	return &AmqpWorkspaceConsumer{
		connection:        connection,
		senderDispatchers: make(map[string]*SenderDispatcher, 1000),
		sendDelivery:      make(chan dto.PeerEnvelope),
	}
}

func (apc *AmqpWorkspaceConsumer) Init() error {
	channel, err := apc.connection.Channel()
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

func (apc *AmqpWorkspaceConsumer) CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionBase, error) {
	apc.senderDispatchersMux.Lock()
	defer apc.senderDispatchersMux.Unlock()

	dispatcher, exist := apc.senderDispatchers[payload.WorkspaceID]
	if !exist {
		dispatcher = NewSenderDispatcher(apc.connection, apc.sendDelivery, senderID)
		apc.senderDispatchers[payload.WorkspaceID] = dispatcher
	}
	queue, err := dispatcher.AddQueue(senderID, payload.WorkspaceID)
	if err != nil {
		return nil, err
	}

	config := config.GetConfig().AMQP
	return &dto.CreateWorkspaceConnectionBase{
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

func (apc *AmqpWorkspaceConsumer) Consume(handler func(payload dto.PeerEnvelope)) {
	go func() {
		for payload := range apc.sendDelivery {
			go handler(payload)
		}
	}()
}
