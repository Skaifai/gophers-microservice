package logger

import (
	"context"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"time"
)

type Publisher struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewPublisher(cfg *config.Config) (*Publisher, error) {
	rmqPort := strconv.Itoa(cfg.RMQ.Port)
	rmqDSN := fmt.Sprintf("amqp://" + cfg.RMQ.Username + ":" + cfg.RMQ.Password +
		"@" + cfg.RMQ.Host + ":" + rmqPort)
	log.Println("NewPublisher: RMQ DSN = " + rmqDSN)
	conn, err := amqp.Dial(rmqDSN)
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
