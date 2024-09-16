package consumer

import (
	"go-kafka/internal/config"
	"go-kafka/internal/lib/logger/sl"
	"log/slog"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ConsumerService struct {
	log      *slog.Logger
	cfg      config.Consumer
	consumer *kafka.Consumer
	stop     chan bool
}

func NewConsumerService(log *slog.Logger, cfg config.Consumer) (*ConsumerService, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaBrokers,
		"group.id":          cfg.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Error("failed to create consumer", sl.Err(err))
		return nil, err
	}

	return &ConsumerService{
		log:      log,
		cfg:      cfg,
		consumer: consumer,
		stop:     make(chan bool),
	}, nil
}

func (c *ConsumerService) StopConsumer() {
	c.log.Info("Stopping consumer")
	c.stop <- true
	c.consumer.Close()
}

func (c *ConsumerService) ConsumeMessages() {
	const op = "services.consumer.ConsumeMessages"

	c.log = slog.With(
		slog.String("op", op),
	)
	err := c.consumer.Assign([]kafka.TopicPartition{
		{
			Topic:     &c.cfg.MessageTopic,
			Partition: 0,
		},
		{
			Topic:     &c.cfg.MessageTopic,
			Partition: 1,
		},
	})
	if err != nil {
		c.log.Error("failed connect to topic partition", sl.Err(err))
		os.Exit(1)
	}

	for {
		select {
		case <-c.stop:
			c.log.Info("consumer stopped")
			return
		default:
			msg, err := c.consumer.ReadMessage(time.Second)
			if err == nil {
				c.log.Info(
					"Consumed message",
					slog.String("topic", *msg.TopicPartition.Topic),
					slog.Int("partition", int(msg.TopicPartition.Partition)),
					slog.Any("offset", msg.TopicPartition.Offset),
					slog.String("value", string(msg.Value)),
				)
			} else if !err.(kafka.Error).IsTimeout() {
				c.log.Error("consume message error", sl.Err(err))
			}
		}
	}
}
