package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

type Publisher struct {
	rabbit *rabbitmq.Connection
	queue  string
}

func NewPublisher(r *rabbitmq.Connection, queue string) *Publisher {
	return &Publisher{
		rabbit: r,
		queue:  queue,
	}
}

func (p *Publisher) Publish(ctx context.Context, job EmailJob) error {

	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal email job: %w", err)
	}

	err = p.rabbit.Channel.Publish(
		"",      
		p.queue, 
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, 
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
