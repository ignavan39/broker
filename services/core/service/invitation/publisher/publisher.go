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
	err             chan error
	ping            chan int
}

func NewPublisher(connection *amqp.Connection) *Publisher {
	return &Publisher{
		connection:  connection,
		connections: make(map[string]*InvitationQueue),
		err:         make(chan error),
		ping:        make(chan int),
	}
}

func (p *Publisher) Ping() chan int {
	return p.ping
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

func (p *Publisher) SetLastUpdateTimeByUserId(userID string, time time.Time) {
	p.connectionsLock.Lock()
	defer p.connectionsLock.Unlock()

	queue, ok := p.connections[userID]
	if !ok {
		return
	}

	queue.lastModified = time
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

func (p *Publisher) RemoveDeadQueues(expireTime time.Time) ([]string, error) {
	queueNames := make([]string, 0)
	for id, queue := range p.connections {
		if expireTime.Add(time.Duration(-30) * time.Second).After(queue.LastModified()) {
			if err := queue.Remove(); err != nil {
				return nil, err
			}
			queueNames = append(queueNames, queue.Name())
			delete(p.connections, id)
		}
	}

	return queueNames, nil
}

func (p *Publisher) GetQueueByUser(userID string) (*InvitationQueue) {
	queue, ok := p.connections[userID]

	if !ok {
		return nil
	}

	return queue
}
 