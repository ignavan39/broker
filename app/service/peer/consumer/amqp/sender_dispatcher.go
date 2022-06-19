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
	conn              *amqp.Connection
	senderQueueMap    map[string]*SenderWorkspaceQueue
	senderQueueMapMux sync.RWMutex
	delivery          chan dto.PeerEnvelope
	senderID          string
}

func NewSenderDispatcher(conn *amqp.Connection, out chan dto.PeerEnvelope, senderID string) *SenderDispatcher {
	return &SenderDispatcher{
		senderQueueMap: make(map[string]*SenderWorkspaceQueue, 1000),
		conn:           conn,
		delivery:       out,
		senderID:       senderID,
	}
}

func (d *SenderDispatcher) AddQueue(senderID string, workspaceID string) (*SenderWorkspaceQueue, error) {
	d.senderQueueMapMux.Lock()
	defer d.senderQueueMapMux.Unlock()

	queue, err := NewSenderWorkspaceQueue(senderID, workspaceID, d.conn)

	if err != nil {
		blogger.Infof("[SenderDispatcher][sender :%s] failed add queue %s", senderID, err.Error())
		return nil, err
	}

	exQueue, exist := d.senderQueueMap[senderID]
	if exist {
		return exQueue, nil
	}

	go queue.Run(d.delivery)
	d.senderQueueMap[queue.Name()] = queue

	return queue, nil
}

func (d *SenderDispatcher) GetSenderID() string {
	return d.senderID
}
