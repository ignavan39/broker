package amqp

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/pkg/utils"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type ReceiverQueue struct {
	receiverID  string
	workspaceID string
	meta        dto.Meta
	stop        chan int
	queue       amqp.Queue
	channel     *amqp.Channel
}

func NewReceiverQueue(
	receiverID string,
	workspaceID string,
	connection *amqp.Connection,
) (*ReceiverQueue, error) {
	amqpChannel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	queueName := fmt.Sprintf(
		"%s.%s",
		getPrefix(),
		utils.CryptString(fmt.Sprintf("%s%s", workspaceID, receiverID), config.GetConfig().AMQP.QueueHashSalt))
	exchangeName := fmt.Sprintf("%s-%s", getExchangePrefix(), receiverID)

	queue, err := amqpChannel.QueueDeclare(queueName,
		false,
		true,
		false,
		false,
		nil,
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

	return &ReceiverQueue{
		queue:       queue,
		receiverID:  receiverID,
		workspaceID: workspaceID,
		channel:     amqpChannel,
		stop:        make(chan int),
		meta: dto.Meta{
			QueueName:    queueName,
			ExchangeName: exchangeName,
		},
	}, nil
}

func (rq *ReceiverQueue) Publish(payload dto.PeerEnvelope) error {
	buffer, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[ReceiverQueue] error serialize payload %s, [Error] %s", rq.meta.QueueName, err.Error())
	}

	return rq.channel.Publish(
		rq.meta.ExchangeName,
		rq.meta.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        buffer,
		})
}

func (rq *ReceiverQueue) Name() string {
	return rq.meta.QueueName
}

func (p *ReceiverQueue) ExchangeName() string {
	return p.meta.ExchangeName
}

func getExchangePrefix() string {
	return "workspace"
}

func getPrefix() string {
	return "workspace.publisher"
}
