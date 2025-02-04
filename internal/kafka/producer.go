package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)


type Producer struct {
	Writer *kafka.Writer
}


func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
	})
	return &Producer{
		Writer: writer,
	}
}


func (p *Producer) SendMessage(key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}
	return p.Writer.WriteMessages(context.Background(), msg)
}


func (p *Producer) Close() error {
	return p.Writer.Close()
}
