package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Connection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func New(url string) (*Connection, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Connection{
		Conn:    conn,
		Channel: ch,
	}, nil
}
