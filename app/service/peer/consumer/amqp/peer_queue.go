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

type PeerQueue struct {
	recipientID string
	senderID    string
	meta        dto.Meta
	stop        chan int
	queue       amqp.Queue
	channel     *amqp.Channel
}

func NewPeerQueue(recipientID string, senderID string, workspaceID string, conn *amqp.Connection) (*PeerQueue, error) {
	amqpChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queueName := utils.CryptString(fmt.Sprintf("%s%s", workspaceID, recipientID), config.GetConfig().AMQP.QueueHashSalt)
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

	return &PeerQueue{
		queue:       queue,
		recipientID: recipientID,
		senderID:    senderID,
		channel:     amqpChannel,
		stop:        make(chan int),
		meta: dto.Meta{
			QueueName:    queueName,
			ExchangeName: exchangeName,
			ReportKey:    queueName,
		},
	}, nil
}

func (q *PeerQueue) Run(out chan dto.PeerEnvelope) error {
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
	for delivery := range deliveries{
		var payload dto.PeerEnvelope
		err := json.Unmarshal(delivery.Body, &payload)
		if err != nil {
			blogger.Errorf("[ClientQueue][Queue :%s] failed decode", q.meta.QueueName)
		} else {
			out <- payload
		}
	}
	return nil
}
func getExchangePrefix() string {
	return "peer"
}
