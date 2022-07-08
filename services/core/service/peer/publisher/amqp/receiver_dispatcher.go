package amqp

import (
	"broker/pkg/logger"
	"errors"
	"sync"

	"github.com/streadway/amqp"
)

var QueueNameAlreadyExist = errors.New("queue by peer already exist")

type ReceiverDispatcher struct {
	connection          *amqp.Connection
	receiveQueueMap     map[string]*ReceiverQueue
	receiveQueueMapLock sync.RWMutex
}

func NewReceiverDispatcher(connection *amqp.Connection) *ReceiverDispatcher {
	return &ReceiverDispatcher{
		receiveQueueMap: make(map[string]*ReceiverQueue, 1000),
		connection:      connection,
	}
}

func (d *ReceiverDispatcher) AddQueue(receiveID string, workspaceID string) (*ReceiverQueue, error) {
	d.receiveQueueMapLock.Lock()
	defer d.receiveQueueMapLock.Unlock()

	existingQueue, exist := d.receiveQueueMap[receiveID]
	if exist {
		return existingQueue, nil
	} else {
		queue, err := NewReceiverQueue(receiveID, workspaceID, d.connection)

		if err != nil {
			logger.Logger.Infof("[ReceiverDispatcher][receive :%s] failed add queue %s", receiveID, err.Error())
			return nil, err
		}

		d.receiveQueueMap[queue.Name()] = queue

		return queue, nil
	}
}

func (d *ReceiverDispatcher) GetQueue(receiveID string) (*ReceiverQueue, bool) {
	d.receiveQueueMapLock.Lock()
	defer d.receiveQueueMapLock.Unlock()

	queue, exist := d.receiveQueueMap[receiveID]

	return queue, exist
}
