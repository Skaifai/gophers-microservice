package logger

import "C"
import (
	"context"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Publisher struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewPublisher() (*Publisher, error) {
	conn, err := amqp.Dial(config.GetEnvironmentVar("RMQ_DSN"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (p *Publisher) SendLog(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.Channel.PublishWithContext(ctx,
		"logs",   // exchange
		"logger", // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Publisher) Close() error {
	err := p.Channel.Close()
	if err != nil {
		return err
	}
	return p.Conn.Close()
}
