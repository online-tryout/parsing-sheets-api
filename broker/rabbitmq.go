package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	Channel *amqp091.Channel
}

type Message struct {
	ProcessID string `json:"processId"`
	EndedAt   string `json:"endedAt"`
	Price     string `json:"price"`
	StartedAt string `json:"startedAt"`
	Status    string `json:"status"`
	Title     string `json:"title"`
	URL       string `json:"url"`
}

func NewRabbitMq(source string) (*RabbitMq, error) {
	con, err := amqp091.Dial(source)
	if err != nil {
		return nil, err
	}

	rcon, err := con.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMq{
		Channel: rcon,
	}, nil
}

func (rmq *RabbitMq) PublishEvent(queue string, msg []byte) error {
	q, err := rmq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publishedMsg := amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         msg,
	}

	err = rmq.Channel.PublishWithContext(ctx, "", q.Name, false, false, publishedMsg)
	if err != nil {
		return err
	}

	return nil
}

func (rmq *RabbitMq) ConsumeEvent(queue string) error {
	q, err := rmq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := rmq.Channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		err := handleMessage(msg.Body)
		if err != nil {
			msg.Nack(false, true)
			return err
		}
		msg.Ack(false)
	}

	return nil
}

func handleMessage(body []byte) error {
	var msg Message

	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	fmt.Printf("Received message:\nProcessID: %s\nEndedAt: %s\nPrice: %s\nStartedAt: %s\nStatus: %s\nTitle: %s\nURL: %s\n",
		msg.ProcessID, msg.EndedAt, msg.Price, msg.StartedAt, msg.Status, msg.Title, msg.URL)
	return nil
}