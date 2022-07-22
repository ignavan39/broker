package consumer

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/pkg/logger"
	"broker/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type SenderQueue struct {
	senderID    string
	workspaceID string
	meta        dto.Meta
	stop        chan int
	queue       amqp.Queue
	channel     *amqp.Channel
}

func NewSenderQueue(
	senderID string,
	workspaceID string,
	connection *amqp.Connection,
) (*SenderQueue, error) {
	amqpChannel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	queueName := fmt.Sprintf(
		"%s.%s",
		getPrefix(),
		utils.CryptString(fmt.Sprintf("%s%s", workspaceID, senderID), config.GetConfig().AMQP.QueueHashSalt))
	exchangeName := fmt.Sprintf("%s-%s", getExchangePrefix(), senderID)

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

	return &SenderQueue{
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

func (q *SenderQueue) Run(out chan dto.PeerEnvelope) error {
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
			logger.Logger.Errorf("[WorkspaceQueue][Queue :%s] failed decode", q.meta.QueueName)
			delivery.Nack(false, false)
			continue
		} else {
			logger.Logger.Infof("[WorkspaceQueue][Queue :%s] receive message %v", payload)
			delivery.Ack(false)
			out <- payload
		}
	}
	return nil
}

func (p *SenderQueue) Name() string {
	return p.meta.QueueName
}

func (p *SenderQueue) ExchangeName() string {
	return p.meta.ExchangeName
}

func getPrefix() string {
	return "workspace.consumer"
}

func getExchangePrefix() string {
	return "workspace"
}

func getDeadLetterExchangeName() string {
	return "workspace.dead-letter"
}

func getDeadLetterQueueName() string {
	return "workspace.dead-letter.queue"
}
