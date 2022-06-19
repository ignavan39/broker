package amqp

import (
	"broker/app/dto"
	"errors"
	"sync"

	"github.com/streadway/amqp"

	blogger "github.com/sirupsen/logrus"
)

var QueueNameAlreadyExist = errors.New("queue by peer already exist")

type SenderDispatcher struct {
	connection         *amqp.Connection
	senderQueueMap     map[string]*SenderQueue
	senderQueueMapLock sync.RWMutex
	delivery           chan dto.PeerEnvelope
}

func NewSenderDispatcher(connection *amqp.Connection, out chan dto.PeerEnvelope) *SenderDispatcher {
	return &SenderDispatcher{
		senderQueueMap: make(map[string]*SenderQueue, 1000),
		connection:     connection,
		delivery:       out,
	}
}

func (d *SenderDispatcher) AddQueue(senderID string, workspaceID string) (*SenderQueue, error) {
	d.senderQueueMapLock.Lock()
	defer d.senderQueueMapLock.Unlock()

	existingQueue, exist := d.senderQueueMap[senderID]
	if exist {
		return existingQueue, nil
	} else {
		queue, err := NewSenderQueue(senderID, workspaceID, d.connection)

		if err != nil {
			blogger.Infof("[SenderDispatcher][sender :%s] failed add queue %s", senderID, err.Error())
			return nil, err
		}

		go queue.Run(d.delivery)
		d.senderQueueMap[queue.Name()] = queue

		return queue, nil
	}
}
