package publisher

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"context"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	connection      *amqp.Connection
	connections     map[string]*InvitationQueue
	connectionsLock sync.RWMutex
}

func NewPublisher(connection *amqp.Connection) *Publisher {
	return &Publisher{
		connection:  connection,
		connections: make(map[string]*InvitationQueue),
	}
}

func (p *Publisher) CreateConnection(ctx context.Context, userID string) (*dto.ConnectInvitationResponse, error) {
	p.connectionsLock.Lock()
	defer p.connectionsLock.Unlock()

	queue, err := NewInvitationQueue(userID, *p.connection)

	if err != nil {
		return nil, err
	}

	p.connections[userID] = queue

	config := config.GetConfig().AMQP
	return &dto.ConnectInvitationResponse{
		Consume: dto.CreateConnectionBase{
			Host:     config.Host,
			Port:     config.Port,
			User:     config.ExternalUser,
			Password: config.ExternalPassword,
			Vhost:    config.Vhost,
			Meta: dto.Meta{
				QueueName:    queue.Name(),
				ExchangeName: queue.ExchangeName(),
			},
		},
	}, nil
}

func (p *Publisher) Publish(userID string, invitation models.Invitation) error {
	p.connectionsLock.Lock()
	defer p.connectionsLock.Unlock()

	queue, ok := p.connections[userID]

	if !ok {
		return nil
	}

	err := queue.Publish(invitation)

	if err != nil {
		return err
	}

	return nil
}

func (p *Publisher) DeleteQueue(invitationQueue *InvitationQueue) {
	p.connectionsLock.Lock()
	defer p.connectionsLock.Unlock()

	for userID, queue := range p.connections {
		if queue == invitationQueue {
			delete(p.connections, userID)
		}
	}
}

func (p *Publisher) GetExpiredQueues(expireTime time.Time) []*InvitationQueue {
	queues := make([]*InvitationQueue, 0)

	for _, queue := range p.connections {
		if expireTime.Add(time.Duration(-30) * time.Second).After(queue.LastModified()) {
			queues = append(queues, queue)
		}
	}

	return queues
}
