package model

import (
	"encoding/json"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp091.Connection
	ch *amqp091.Channel
	queue amqp091.Queue
}

type Params struct {
	To string `json:"to"`
	Subject string `json:"subject"`
	TemplateId int `json:"templateId"`
	Params map[string]interface{} `json:"params"`
}

func NewRabbitMQ(queueName string) *RabbitMQ {
	rmq := &RabbitMQ{}

	conn, err := amqp091.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Fatalf("Falha ao conectar ao servidor AMQP: %v", err)
	}
	rmq.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Falha ao abrir canal: %v", err)
	}
	rmq.ch = ch

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false, 
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Falha ao declarar fila: %v", err)
	}
	rmq.queue = queue

	e := rmq.ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		log.Fatalf("Falha ao declarar exchange: %v", e)
	}

	err = ch.QueueBind(
		rmq.queue.Name,
		"",
		"logs",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Falha ao vincular fila à exchange: %v", err)
	}

	return rmq
}

func (r *RabbitMQ) Close() {
	if r.ch != nil {
		r.ch.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQ) Publish(to string, logs Log) error {
	body := Params{
		To: to,
		Subject: "Service unavailable",
		TemplateId: 3,
		Params: map[string]interface{}{
			"url": logs.URL,
			"status": logs.Status,
			"statusCode": logs.StatusCode,
			"message": logs.Message,
			"elapsedTime": logs.ElapsedTime,
			"dateTime": logs.DateTime,
		},
	}

	jso, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Falha ao serializar log: %v", err)
	}

	err = r.ch.Publish(
		"logs",
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body: jso,
		},
	)
	if err != nil {
		log.Fatalf("Falha ao publicar mensagem: %v", err)
	}

	return nil
}