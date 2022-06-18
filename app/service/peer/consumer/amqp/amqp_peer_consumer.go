package amqp

import (
	"broker/app/config"
	"broker/app/dto"
	"context"
	"sync"

	"github.com/streadway/amqp"
)

type AmqpPeerConsumer struct {
	connection              *amqp.Connection
	receiveDelivery         chan dto.PeerEnvelope
	recipientDispatchers    map[string]*RecipientDispatcher
	recipientDispatchersMux sync.RWMutex
}

func NewAmqpPeerConsumer(connection *amqp.Connection) *AmqpPeerConsumer {
	return &AmqpPeerConsumer{
		connection:           connection,
		recipientDispatchers: make(map[string]*RecipientDispatcher, 1000),
		receiveDelivery:      make(chan dto.PeerEnvelope),
	}
}

func (apc *AmqpPeerConsumer) CreateConnection(ctx context.Context,recipientID string, payload dto.CreatePeerConnectionPayload) (*dto.CreatePeerConnectionResponse, error) {
	apc.recipientDispatchersMux.Lock()
	defer apc.recipientDispatchersMux.Unlock()

	dispatcher,exist := apc.recipientDispatchers[recipientID]
	if !exist {
		dispatcher := NewRecipientDispatcher(apc.connection,apc.receiveDelivery,recipientID)
		apc.recipientDispatchers[recipientID] = dispatcher
	}

	queue,err := dispatcher.AddQueue(recipientID,payload.PeerId)
	if err != nil {
		return nil,err
	}

	config := config.GetConfig().AMQP
	return &dto.CreatePeerConnectionResponse{
		Host: config.Host,
		Port: config.Port,
		User: config.User,
		Password: config.Pass,
		Vhost: config.Vhost,
		Meta: dto.Meta{
			QueueName: queue.Name(),
			ExchangeName: queue.ExchangeName(),
		},
	},nil

}
