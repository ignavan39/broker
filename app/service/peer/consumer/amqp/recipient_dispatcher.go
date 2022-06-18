package amqp

import (
	"broker/app/dto"
	"errors"
	"sync"

	"github.com/streadway/amqp"

	blogger "github.com/sirupsen/logrus"
)

var QueueNameAlreadyExist = errors.New("queue by peer already exist")

type RecipientDispatcher struct {
	conn            *amqp.Connection
	RecipientPeerQueueMap    map[string]*RecipientPeerQueue
	RecipientPeerQueueMapMux sync.RWMutex
	delivery        chan dto.PeerEnvelope
	recipientID     string
}

func NewRecipientDispatcher(conn *amqp.Connection, out chan dto.PeerEnvelope, recipientID string) *RecipientDispatcher {
	return &RecipientDispatcher{
		RecipientPeerQueueMap: make(map[string]*RecipientPeerQueue, 1000),
		conn:         conn,
		delivery:     out,
		recipientID:  recipientID,
	}
}

func (d *RecipientDispatcher) AddQueue(recipientID string, peerID string) (*RecipientPeerQueue, error) {
	d.RecipientPeerQueueMapMux.Lock()
	defer d.RecipientPeerQueueMapMux.Unlock()

	queue, err := NewRecipientPeerQueue(recipientID, peerID, d.conn)

	if err != nil {
		blogger.Infof("[RecipientDispatcher][recipient :%s] failed add queue %s", recipientID,err.Error())
		return nil,err
	}

	exQueue, exist := d.RecipientPeerQueueMap[recipientID]
	if exist {
		return exQueue,nil
	}

	go queue.Run(d.delivery)
	d.RecipientPeerQueueMap[queue.Name()] = queue

	return queue,nil
}

func (d *RecipientDispatcher) GetRecipientID () string {
	return d.recipientID
}