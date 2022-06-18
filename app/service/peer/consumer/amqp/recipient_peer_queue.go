package amqp

import (
	"broker/app/config"
	"broker/app/dto"
	"broker/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	blogger "github.com/sirupsen/logrus"
)

type RecipientPeerQueue struct {
	recipientID string
	peerID    string
	meta        dto.Meta
	stop        chan int
	queue       amqp.Queue
	channel     *amqp.Channel
}

func NewRecipientPeerQueue(recipientID string, peerID string, conn *amqp.Connection) (*RecipientPeerQueue, error) {
	amqpChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queueName := utils.CryptString(fmt.Sprintf("%s%s", peerID, recipientID), config.GetConfig().AMQP.QueueHashSalt)
	exchangeName := fmt.Sprintf("%s_%s", getExchangePrefix(), recipientID)
	queue, err := amqpChannel.QueueDeclare(queueName, false, true, true, false, nil)
	if err != nil {
		return nil, err
	}

	err = amqpChannel.ExchangeDeclare(exchangeName, "topic", false, true, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = amqpChannel.QueueBind(queueName, queueName, exchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	return &RecipientPeerQueue{
		queue:       queue,
		recipientID: recipientID,
		peerID: peerID,
		channel:     amqpChannel,
		stop:        make(chan int),
		meta: dto.Meta{
			QueueName:    queueName,
			ExchangeName: exchangeName,
		},
	}, nil
}

func (q *RecipientPeerQueue) Run(out chan dto.PeerEnvelope) error {
	deliveries, err := q.channel.Consume(
		q.meta.QueueName,
		fmt.Sprintf("consumer-%s", q.meta.QueueName),
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		var payload dto.PeerEnvelope
		err := json.Unmarshal(delivery.Body, &payload)
		if err != nil {
			blogger.Errorf("[PeerQueue][Queue :%s] failed decode", q.meta.QueueName)
			continue
		} else {
			blogger.Infof("[PeerQueue][Queue :%s] receive message %v", payload)
			out <- payload
		}
	}
	return nil
}


func (p *RecipientPeerQueue) Name() string {
	return p.meta.QueueName
}

func (p *RecipientPeerQueue) ExchangeName() string {
	return p.meta.ExchangeName
}

func getExchangePrefix() string {
	return "peer"
}
