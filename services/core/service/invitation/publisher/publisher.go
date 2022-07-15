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
	connections     map[string]map[string]*InvitationQueue
	connectionsLock sync.RWMutex
	err             chan error
	ping            chan int
}

func NewPublisher(connection *amqp.Connection) *Publisher {
	return &Publisher{
		connection:  connection,
		connections: make(map[string]map[string]*InvitationQueue),
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

	_, ok := p.connections[userID]

	if !ok {
		qs := make(map[string]*InvitationQueue)
		p.connections[userID] = qs
	}

	p.connections[userID][queue.meta.QueueName] = queue

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

	queues, ok := p.connections[userID]
	if !ok {
		return
	}

	for _, queue := range queues {
		queues[queue.meta.QueueName].lastModified = time
	}
}

func (p *Publisher) Publish(userID string, invitation models.Invitation) error {
	p.connectionsLock.Lock()
	defer p.connectionsLock.Unlock()

	queues, ok := p.connections[userID]

	if !ok {
		return nil
	}

	for _, queue := range queues {
		err := queues[queue.meta.QueueName].Publish(invitation)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Publisher) RemoveDeadQueues(expireTime time.Time) ([]string, error) {
	queueNames := make([]string, 0)
	for _, queues := range p.connections {
		for id, queue := range queues {
			if expireTime.Add(time.Duration(-30) * time.Second).After(queue.LastModified()) {
				if err := queue.Remove(); err != nil {
					return nil, err
				}
				queueNames = append(queueNames, queue.Name())
				delete(p.connections, id)
			}
		}
	}

	return queueNames, nil
}

func (p *Publisher) GetQueueByUser(userID string, queueName string) *InvitationQueue {
	queues, ok := p.connections[userID]

	if !ok {
		return nil
	}

	return queues[queueName]
}
