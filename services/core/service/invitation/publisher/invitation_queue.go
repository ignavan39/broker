package publisher

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"broker/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type InvitationQueue struct {
	userID       string
	meta         dto.Meta
	queue        amqp.Queue
	channel      *amqp.Channel
	lastModified time.Time
}

func NewInvitationQueue(userID string, connection amqp.Connection) (*InvitationQueue, error) {
	amqpChannel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	queueName := fmt.Sprintf(
		"%s.%s",
		getPrefix(),
		utils.CryptString(userID, config.GetConfig().AMQP.QueueHashSalt))

	exchangeName := getExchangePrefix()

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
		true,
		false,
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

	return &InvitationQueue{
		userID: userID,
		meta: dto.Meta{
			QueueName:    queueName,
			ExchangeName: exchangeName,
		},
		queue:   queue,
		channel: amqpChannel,
	}, nil
}

func (q *InvitationQueue) Publish(invitation models.Invitation) error {
	buffer, err := json.Marshal(invitation)
	if err != nil {
		return fmt.Errorf("[InvitationQueue] error serialize payload %s, [Error] %s", q.meta.QueueName, err.Error())
	}

	return q.channel.Publish(
		q.meta.ExchangeName,
		q.meta.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        buffer,
		})
}

func (q *InvitationQueue) Remove() error {
	defer q.channel.Close()
	_, err := q.channel.QueueDelete(q.Name(), false, false, false)

	if err != nil {
		return err
	}

	return nil
}

func getExchangePrefix() string {
	return "invitation"
}

func getPrefix() string {
	return "invitation.publisher"
}

func (q *InvitationQueue) Name() string {
	return q.meta.QueueName
}

func (q *InvitationQueue) ExchangeName() string {
	return q.meta.ExchangeName
}

func (q *InvitationQueue) LastModified() time.Time {
	return q.lastModified
}
