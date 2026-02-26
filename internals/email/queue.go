package email

import (
	"fmt"

	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

func SetupEmailQueue(r *rabbitmq.Connection,queuename string)error{
	_,err := r.Channel.QueueDeclare(
		queuename,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue;%w",err)
	}
	return nil
}