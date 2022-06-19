package amqp

import (
	"broker/app/config"
	"broker/app/dto"
	"context"
	"sync"

	"github.com/streadway/amqp"
)

type Publisher struct {
	connection             *amqp.Connection
	receiveDispatchers     map[string]*ReceiverDispatcher
	receiveDispatchersLock sync.RWMutex
}

func NewPublisher(connection *amqp.Connection) *Publisher {
	return &Publisher{
		connection:         connection,
		receiveDispatchers: make(map[string]*ReceiverDispatcher, 1000),
	}
}

func (p *Publisher) CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionBase, error) {
	p.receiveDispatchersLock.Lock()
	defer p.receiveDispatchersLock.Unlock()

	dispatcher, exist := p.receiveDispatchers[payload.WorkspaceID]
	if !exist {
		dispatcher = NewReceiverDispatcher(p.connection)
		p.receiveDispatchers[payload.WorkspaceID] = dispatcher
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

func (p *Publisher) Publish(workspaceID string, payload dto.PeerEnvelope) error {
	p.receiveDispatchersLock.Lock()
	defer p.receiveDispatchersLock.Unlock()

	dispatcher, exist := p.receiveDispatchers[workspaceID]
	if !exist {
		return nil
	}

	queue, exist := dispatcher.GetQueue(payload.FromId)
	if exist {
		return queue.Publish(payload)
	}

	return nil
}
