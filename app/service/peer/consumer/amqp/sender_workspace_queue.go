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

type SenderWorkspaceQueue struct {
	senderID    string
	workspaceID string
	meta        dto.Meta
	stop        chan int
	queue       amqp.Queue
	channel     *amqp.Channel
}

func NewSenderWorkspaceQueue(
	senderID string,
	workspaceID string,
	conn *amqp.Connection,
) (*SenderWorkspaceQueue, error) {
	amqpChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queueName := fmt.Sprintf(
		"%s.%s",
		getPrefix(),
		utils.CryptString(fmt.Sprintf("%s%s", workspaceID, senderID), config.GetConfig().AMQP.QueueHashSalt))
	exchangeName := fmt.Sprintf("%s-%s", getPrefix(), senderID)

	queue, err := amqpChannel.QueueDeclare(queueName,
		false,
		true,
		false,
		false,
		amqp.Table{"x-dead-letter-exchange": getDeadLetterExchangeName()},
	)
	if err != nil {
		return nil, err
	}

	if err = amqpChannel.ExchangeDeclare(
		exchangeName,
		"topic",
		false,
		true,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	if err = amqpChannel.QueueBind(
		queueName,
		queueName,
		exchangeName,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return &SenderWorkspaceQueue{
		queue:       queue,
		senderID:    senderID,
		workspaceID: workspaceID,
		channel:     amqpChannel,
		stop:        make(chan int),
		meta: dto.Meta{
			QueueName:    queueName,
			ExchangeName: exchangeName,
		},
	}, nil
}

func (q *SenderWorkspaceQueue) Run(out chan dto.PeerEnvelope) error {
	deliveries, err := q.channel.Consume(
		q.meta.QueueName,
		fmt.Sprintf("%s.consumer-%s", getPrefix(), q.meta.QueueName),
		false,
		false,
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
			blogger.Errorf("[WorkspaceQueue][Queue :%s] failed decode", q.meta.QueueName)
			delivery.Nack(false, false)
			continue
		} else {
			blogger.Infof("[WorkspaceQueue][Queue :%s] receive message %v", payload)
			delivery.Ack(false)
			out <- payload
		}
	}
	return nil
}

func (p *SenderWorkspaceQueue) Name() string {
	return p.meta.QueueName
}

func (p *SenderWorkspaceQueue) ExchangeName() string {
	return p.meta.ExchangeName
}

func getPrefix() string {
	return "workspace"
}

func getDeadLetterExchangeName() string {
	return "workspace.dead-letter"
}

func getDeadLetterQueueName() string {
	return "workspace.dead-letter.queue"
}
