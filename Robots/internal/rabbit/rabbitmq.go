package rabbit

import (
	"RobotService/internal/entities"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "robots"
	exchangeType = "direct"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(amqpURL string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn, channel: ch}, nil
}

func (p *Publisher) Close() {
	_ = p.channel.Close()
	_ = p.conn.Close()
}

func (p *Publisher) Publish(robot *entities.Robot, routingKey string) error {
	body, _ := json.Marshal(robot)

	log.Printf("Отправка в рэббит по routing key: %s", routingKey)
	return p.channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) PublishWithText(message string, routingKey string) error {
	log.Printf("Отправка сообщения в рэббит по routing key: %s", routingKey)
	return p.channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}
